package blockchain

import (
	"bytes"

	"github.com/the-code-innovator/go-blockchain/wallet"
)

// TxInput structure for Input for the BlockChain
type TxInput struct {
	ID        []byte
	Out       int
	Signature []byte
	PublicKey []byte
}

// TxOutput structure for Output for the BlockChainw
type TxOutput struct {
	Value         int
	PublicKeyHash []byte
}

// NewTxOutput to create a new transaction output for the new transaction that is created by every spend
func NewTxOutput(value int, address string) *TxOutput {
	txOut := &TxOutput{value, nil}
	txOut.Lock([]byte(address))
	return txOut
}

// Lock to lock the transaction from spending without authorisation
func (out *TxOutput) Lock(address []byte) {
	publicKeyHash := wallet.Base58Decode(address)
	publicKeyHash = publicKeyHash[1 : len(publicKeyHash)-4]
	// out.PublicKey = publicKeyHash
	out.PublicKeyHash = publicKeyHash
}

// UsesKey to check for unlocking
func (in *TxInput) UsesKey(publicKeyHash []byte) bool {
	lockingHash := wallet.PublicKeyHash(in.PublicKey)
	return bytes.Compare(lockingHash, publicKeyHash) == 0
}

// IsLockedWithKey to verify that the transaction is locked with only the users public key
func (out *TxOutput) IsLockedWithKey(publicKeyHash []byte) bool {
	return bytes.Compare(out.PublicKeyHash, publicKeyHash) == 0
}
