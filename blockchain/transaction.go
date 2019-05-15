package blockchain

import (
	"encoding/hex"
	"fmt"
	"log"
)

// CoinBaseTx for the coin base transaction
func CoinBaseTx(to, data string) *Transaction {
	if data == "" {
		data = fmt.Sprintf("COINS TO %s\n", to)
	}
	txin := TxInput{[]byte{}, -1, data}
	txout := TxOutput{100, to}
	tx := Transaction{nil, []TxInput{txin}, []TxOutput{txout}}
	tx.SetID()
	return &tx
}

// IsCoinBase to check for CoinBase Transaction
func (tx *Transaction) IsCoinBase() bool {
	return len(tx.Inputs) == 1 && len(tx.Inputs[0].ID) == 0 && tx.Inputs[0].Output == -1
}

// NewTransaction for creating a new Transaction in the BlockChain
func NewTransaction(from, to string, amount int, blockchain *BlockChain) *Transaction {
	var inputs []TxInput
	var outputs []TxOutput
	accumulator, validOutputs := blockchain.FindSpendableOutputs(from, amount)
	if accumulator < amount {
		log.Panic("ERROR: NOT ENOUGH FUNDS !")
	}
	for txID, outputs := range validOutputs {
		txIDString, err := hex.DecodeString(txID)
		Handle(err)
		for _, output := range outputs {
			input := TxInput{txIDString, output, from}
			inputs = append(inputs, input)
		}
	}
	outputs = append(outputs, TxOutput{amount, to})
	if accumulator > amount {
		outputs = append(outputs, TxOutput{accumulator - amount, from})
	}
	tx := Transaction{nil, inputs, outputs}
	tx.SetID()
	return &tx
}
