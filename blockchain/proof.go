package blockchain

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"math"
	"math/big"
)

// Take Data from the Block
// Create a Counter (nonce) which starts @ 0
// Create a Hash of the Data + Counter
// Check the Hash to see if it meets a Set of Requirements
// Signs the Block, if meets the Set of Requirements
// Create a Block if not.
// Requirements:
// First few Bytes must contain 0's.

// Difficulty Parameter
const Difficulty = 18

// ProofOfWork structure for the proof of Work
type ProofOfWork struct {
	Block  *Block
	Target *big.Int
}

// NewProof for the new ProofOfWork for the Block
func NewProof(block *Block) *ProofOfWork {
	target := big.NewInt(1)
	target.Lsh(target, uint(256-Difficulty))
	proofOfWork := &ProofOfWork{block, target}
	return proofOfWork
}

// InitData for initializing the data in the Block
func (proofOfWork *ProofOfWork) InitData(nonce int) []byte {
	data := bytes.Join(
		[][]byte{
			proofOfWork.Block.PreviousHash,
			proofOfWork.Block.HashTransactions(),
			ToHex(int64(nonce)),
			ToHex(int64(Difficulty)),
		},
		[]byte{},
	)
	return data
}

// Run the computation for the BlockChain
func (proofOfWork *ProofOfWork) Run() (int, []byte) {
	var intHash big.Int
	var hash [32]byte
	nonce := 0
	for nonce < math.MaxInt64 {
		data := proofOfWork.InitData(nonce)
		hash = sha256.Sum256(data)
		fmt.Printf("\r%x", hash)
		intHash.SetBytes(hash[:])
		if intHash.Cmp(proofOfWork.Target) == -1 {
			break
		} else {
			nonce++
		}
	}
	fmt.Println()
	return nonce, hash[:]
}

// Validate for validation of the ProofOfWork
func (proofOfWork *ProofOfWork) Validate() bool {
	var intHash big.Int
	data := proofOfWork.InitData(proofOfWork.Block.Nonce)
	hash := sha256.Sum256(data)
	intHash.SetBytes(hash[:])
	return intHash.Cmp(proofOfWork.Target) == -1
}

// ToHex to convert an integer to Hex
func ToHex(number int64) []byte {
	buffer := new(bytes.Buffer)
	err := binary.Write(buffer, binary.BigEndian, number)
	Handle(err)
	return buffer.Bytes()
}
