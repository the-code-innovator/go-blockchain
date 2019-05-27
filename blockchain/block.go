package blockchain

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"log"
)

// Block structure for the Block type in the blockchain
type Block struct {
	Hash         []byte
	Transactions []*Transaction
	PreviousHash []byte
	Nonce        int
}

// Genesis to create the genesis block in the blockchain
func Genesis(coinbase *Transaction) *Block {
	return CreateBlock([]*Transaction{coinbase}, []byte{})
}

// CreateBlock to create a block in the blockchain
func CreateBlock(txns []*Transaction, previousHash []byte) *Block {
	block := &Block{[]byte{}, txns, previousHash, 0}
	// block.DeriveHash()
	proofOfWork := NewProof(block)
	nonce, hash := proofOfWork.Run()
	block.Nonce = nonce
	block.Hash = hash
	return block
}

// Serialize to serialize the input to BadgerDB
func (block *Block) Serialize() []byte {
	var result bytes.Buffer
	encoder := gob.NewEncoder(&result)
	err := encoder.Encode(block)
	PanicHandle(err)
	return result.Bytes()
}

// Deserialize to deserialize the output from BadgerDB
func Deserialize(data []byte) *Block {
	var block Block
	decoder := gob.NewDecoder(bytes.NewReader(data))
	err := decoder.Decode(&block)
	PanicHandle(err)
	return &block
}

// HashTransactions to hash the transactions in the block
func (block *Block) HashTransactions() []byte {
	var txHashes [][]byte
	var txHash [32]byte
	for _, tx := range block.Transactions {
		txHashes = append(txHashes, tx.ID)
	}
	txHash = sha256.Sum256(bytes.Join(txHashes, []byte{}))
	return txHash[:]
}

// PanicHandle to Panic throw the Error
func PanicHandle(err error) {
	if err != nil {
		log.Panic(err)
	}
}
