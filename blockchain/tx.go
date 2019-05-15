package blockchain

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
)

// Transaction structure for the Transaction
type Transaction struct {
	ID      []byte
	Inputs  []TxInput
	Outputs []TxOutput
}

// TxOutput for Output for the BlockChainw
type TxOutput struct {
	Value     int
	PublicKey string
}

// TxInput for Input for the BlockChain
type TxInput struct {
	ID        []byte
	Output    int
	Signature string
}

// SetID for setting the ID for the Transaction
func (tx *Transaction) SetID() {
	var encoded bytes.Buffer
	var hash [32]byte
	encoder := gob.NewEncoder(&encoded)
	err := encoder.Encode(tx)
	Handle(err)
	hash = sha256.Sum256(encoded.Bytes())
	tx.ID = hash[:]
}

// CanUnlock for checking who can unlock the coinbase transaction
func (in *TxInput) CanUnlock(data string) bool {
	return in.Signature == data
}

// CanBeUnlocked to check whether the transaction can be unlocked
func (out *TxOutput) CanBeUnlocked(data string) bool {
	return out.PublicKey == data
}
