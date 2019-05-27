package blockchain

import (
	"bytes"
	"crypto/ecdsa"
	"encoding/hex"
	"errors"
	"fmt"
	"os"
	"runtime"

	"github.com/dgraph-io/badger"
)

// constants used in the blockchain
const (
	dbPath      = "./tmp/blocks"
	dbFile      = "./tmp/blocks/MANIFEST"
	genesisData = "FIRST TRANSACTION FROM GENESIS."
)

// BlockChain structure for the BlockChain type in the blockchain
type BlockChain struct {
	LastHash []byte
	DataBase *badger.DB
}

// ChainIterator structure to iterate the Blocks in badger.DB
type ChainIterator struct {
	CurrentHash []byte
	DataBase    *badger.DB
}

// InitBlockChain to initialize the BlockChain
func InitBlockChain(address string) *BlockChain {
	if badgerDBExists() {
		fmt.Println("BLOCKCHAIN ALREADY EXISTS.")
		runtime.Goexit()
	}
	var lastHash []byte
	options := badger.DefaultOptions
	options.Dir = dbPath
	options.ValueDir = dbPath
	database, err := badger.Open(options)
	PanicHandle(err)
	err = database.Update(func(txn *badger.Txn) error {
		coinBaseTransaction := CoinBaseTx(address, genesisData)
		genesis := Genesis(coinBaseTransaction)
		fmt.Println("GENESIS CREATED.")
		err := txn.Set(genesis.Hash, genesis.Serialize())
		PanicHandle(err)
		err = txn.Set([]byte("lh"), genesis.Hash)
		lastHash = genesis.Hash
		return err
	})
	PanicHandle(err)
	blockchain := BlockChain{lastHash, database}
	return &blockchain
}

// ContinueBlockChain to continue blockchain validation
func ContinueBlockChain(address string) *BlockChain {
	if badgerDBExists() == false {
		fmt.Println("NO EXISTING BLOCKCHAIN FOUND.\nCREATE ONE.")
		runtime.Goexit()
	}
	var lastHash []byte
	options := badger.DefaultOptions
	options.Dir = dbPath
	options.ValueDir = dbPath
	database, err := badger.Open(options)
	PanicHandle(err)
	err = database.Update(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte("lh"))
		PanicHandle(err)
		lastHash, err = item.Value()
		return err
	})
	PanicHandle(err)
	blockchain := BlockChain{lastHash, database}
	return &blockchain
}

// SignTransaction to sign the transaction that is added to a block
func (chain *BlockChain) SignTransaction(tx *Transaction, privateKey ecdsa.PrivateKey) {
	previousTXs := make(map[string]Transaction)
	for _, in := range tx.Inputs {
		previousTX, err := chain.FindTransaction(in.ID)
		PanicHandle(err)
		previousTXs[hex.EncodeToString(previousTX.ID)] = previousTX
	}
	tx.Sign(privateKey, previousTXs)
}

// VerifyTransaction to verify the transactions in a block
func (chain *BlockChain) VerifyTransaction(tx *Transaction) bool {
	previousTXs := make(map[string]Transaction)
	for _, in := range tx.Inputs {
		previousTX, err := chain.FindTransaction(in.ID)
		PanicHandle(err)
		previousTXs[hex.EncodeToString(previousTX.ID)] = previousTX
	}
	return tx.Verify(previousTXs)
}

// FindTransaction to find a transaction by ID in the list of transactions in the blocks
func (chain *BlockChain) FindTransaction(ID []byte) (Transaction, error) {
	iterator := chain.Iterator()
	for {
		block := iterator.Next()
		for _, tx := range block.Transactions {
			if bytes.Compare(tx.ID, ID) == 0 {
				return *tx, nil
			}
		}
		if len(block.PreviousHash) == 0 {
			break
		}
	}
	return Transaction{}, errors.New("transaction doesn't exist")
}

// AddBlock to add a block to the existing BlockChain
func (chain *BlockChain) AddBlock(transactions []*Transaction) {
	var lastHash []byte
	err := chain.DataBase.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte("lh"))
		PanicHandle(err)
		lastHash, err = item.Value()
		return err
	})
	PanicHandle(err)
	newBlock := CreateBlock(transactions, lastHash)
	err = chain.DataBase.Update(func(txn *badger.Txn) error {
		err := txn.Set(newBlock.Hash, newBlock.Serialize())
		PanicHandle(err)
		err = txn.Set([]byte("lh"), newBlock.Hash)
		chain.LastHash = newBlock.Hash
		return err
	})
	PanicHandle(err)
}

// FindUnspentTransactions to find unspent transactions in the blockchain
func (chain *BlockChain) FindUnspentTransactions(publicKeyHash []byte) []Transaction {
	var unSpentTransactions []Transaction
	spentTxns := make(map[string][]int)
	iterator := chain.Iterator()
	for {
		block := iterator.Next()
		for _, tx := range block.Transactions {
			txID := hex.EncodeToString(tx.ID)
		OutputIterate:
			for outID, out := range tx.Outputs {
				if spentTxns[txID] != nil {
					for _, spentOut := range spentTxns[txID] {
						if spentOut == outID {
							continue OutputIterate
						}
					}
				}
				if out.IsLockedWithKey(publicKeyHash) {
					unSpentTransactions = append(unSpentTransactions, *tx)
				}
			}
			if tx.IsCoinBase() == false {
				for _, in := range tx.Inputs {
					if in.UsesKey(publicKeyHash) {
						inTxID := hex.EncodeToString(in.ID)
						spentTxns[inTxID] = append(spentTxns[inTxID], in.Out)
					}
				}
			}
		}
		if len(block.PreviousHash) == 0 {
			break
		}
	}
	return unSpentTransactions
}

// FindSpendableOutputs to find spendable outputs in the BlockChain
func (chain *BlockChain) FindSpendableOutputs(publicKeyHash []byte, amount int) (int, map[string][]int) {
	unSpentOutputs := make(map[string][]int)
	unSpentTransactions := chain.FindUnspentTransactions(publicKeyHash)
	accumulated := 0
Work:
	for _, tx := range unSpentTransactions {
		txID := hex.EncodeToString(tx.ID)
		for outID, output := range tx.Outputs {
			if output.IsLockedWithKey(publicKeyHash) && accumulated < amount {
				accumulated += output.Value
				unSpentOutputs[txID] = append(unSpentOutputs[txID], outID)
				if accumulated >= amount {
					break Work
				}
			}
		}
	}
	return accumulated, unSpentOutputs
}

// FindUnspentTransactionsOutputs to find unspent transaction outputs in the blockchain
func (chain *BlockChain) FindUnspentTransactionsOutputs(publicKeyHash []byte) []TxOutput {
	var unSpentTransactionOutputs []TxOutput
	unSpentTransactions := chain.FindUnspentTransactions(publicKeyHash)
	for _, tx := range unSpentTransactions {
		for _, output := range tx.Outputs {
			if output.IsLockedWithKey(publicKeyHash) {
				unSpentTransactionOutputs = append(unSpentTransactionOutputs, output)
			}
		}
	}
	return unSpentTransactionOutputs
}

// badgerDBExists to check the availability of DataBase
func badgerDBExists() bool {
	if _, err := os.Stat(dbFile); os.IsNotExist(err) {
		return false
	}
	return true
}

// Iterator to initialise a Iterator over badgerDB
func (chain *BlockChain) Iterator() *ChainIterator {
	iterator := &ChainIterator{chain.LastHash, chain.DataBase}
	return iterator
}

// Next to navigate to the next Block in badgerDB
func (iterator *ChainIterator) Next() *Block {
	var block *Block
	err := iterator.DataBase.View(func(txn *badger.Txn) error {
		item, err := txn.Get(iterator.CurrentHash)
		PanicHandle(err)
		encodedBlock, err := item.Value()
		block = Deserialize(encodedBlock)
		return err
	})
	PanicHandle(err)
	iterator.CurrentHash = block.PreviousHash
	return block
}
