package blockchain

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"log"
)

// Block structure for the Block dataType
type Block struct {
	// Data         []byte
	Hash         []byte
	Transactions []*Transaction
	PreviousHash []byte
	Nonce        int
}

// DeriveHash from the PreviousHash of the same BlockChain
// func (block *Block) DeriveHash() {
// 	info := bytes.Join([][]byte{block.Data, block.PreviousHash}, []byte{})
// 	hash := sha256.Sum256(info)
// 	block.Hash = hash[:]
// }

// HashTransactions to hash the transactions list for the block
func (block *Block) HashTransactions() []byte {
	var txHashes [][]byte
	var txHash [32]byte
	for _, tx := range block.Transactions {
		txHashes = append(txHashes, tx.ID)
	}
	txHash = sha256.Sum256(bytes.Join(txHashes, []byte{}))
	return txHash[:]
}

// CreateBlock from the PreviousHash of the same BlockChain
func CreateBlock(txns []*Transaction, previousHash []byte) *Block {
	block := &Block{[]byte{}, txns, previousHash, 0}
	// block.DeriveHash()
	proofOfWork := NewProof(block)
	nonce, hash := proofOfWork.Run()
	block.Nonce = nonce
	block.Hash = hash
	return block
}

// Genesis of the BlockChain
func Genesis(coinbase *Transaction) *Block {
	// return CreateBlock("Genesis", []byte{})
	return CreateBlock([]*Transaction{coinbase}, []byte{})
}

// Serialize for serializing the output for BadgerDB
func (block *Block) Serialize() []byte {
	var result bytes.Buffer
	encoder := gob.NewEncoder(&result)
	err := encoder.Encode(block)
	Handle(err)
	return result.Bytes()
}

// Deserialize from input BadgerDB for Block
func Deserialize(data []byte) *Block {
	var block Block
	decoder := gob.NewDecoder(bytes.NewReader(data))
	err := decoder.Decode(&block)
	Handle(err)
	return &block
}

// Handle to handle the error
func Handle(err error) {
	if err != nil {
		log.Panic(err)
	}
}
