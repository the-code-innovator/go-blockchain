package blockchain

import (
	"bytes"

	"github.com/the-code-innovator/go-blockchain/wallet"
)

// TxOutput for Output for the BlockChainw
type TxOutput struct {
	Value         int
	PublicKeyHash []byte
}

// TxInput for Input for the BlockChain
type TxInput struct {
	ID        []byte
	Output    int
	Signature []byte
	PublicKey []byte
}

// // CanUnlock for checking who can unlock the coinbase transaction
// func (in *TxInput) CanUnlock(data string) bool {
// 	return in.Signature == data
// }

// // CanBeUnlocked to check whether the transaction can be unlocked
// func (out *TxOutput) CanBeUnlocked(data string) bool {
// 	return out.PublicKey == data
// }

// NewTxOutput to create a new transaction output for the new transaction that is created by every spending
func NewTxOutput(value int, address string) *TxOutput {
	txOut := &TxOutput{value, nil}
	txOut.Lock([]byte(address))
	return txOut
}

// UsesKey to check for unlocking
func (in *TxInput) UsesKey(publicKeyHash []byte) bool {
	lockingHash := wallet.PublicKeyHash(in.PublicKey)
	return bytes.Compare(lockingHash, publicKeyHash) == 0
}

// Lock to lock the transaction from spending without authorisation
func (out *TxOutput) Lock(address []byte) {
	publicKeyHash := wallet.Base58Decode(address)
	publicKeyHash = publicKeyHash[1 : len(publicKeyHash)-4]
	// out.PublicKey = publicKeyHash
	out.PublicKeyHash = publicKeyHash
}

// IsLockedWithKey to verify that the transaction is locked with only the users public key
func (out *TxOutput) IsLockedWithKey(publicKeyHash []byte) bool {
	return bytes.Compare(out.PublicKeyHash, publicKeyHash) == 0
}
