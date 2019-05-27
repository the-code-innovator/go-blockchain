package blockchain

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/gob"
	"encoding/hex"
	"fmt"
	"log"
	"math/big"
	"strings"

	"github.com/the-code-innovator/go-blockchain/wallet"
)

// Transaction structure for the Transaction type in the blockchain
type Transaction struct {
	ID      []byte
	Inputs  []TxInput
	Outputs []TxOutput
}

// CoinBaseTx for the coin base transaction
func CoinBaseTx(to, data string) *Transaction {
	if data == "" {
		data = fmt.Sprintf("COINS TO %s\n", to)
	}
	txin := TxInput{[]byte{}, -1, nil, []byte(data)}
	txout := NewTxOutput(100, to)
	tx := Transaction{nil, []TxInput{txin}, []TxOutput{*txout}}
	tx.SetID()
	return &tx
}

// NewTransaction for creating a new Transaction in the BlockChain
func NewTransaction(from, to string, amount int, blockchain *BlockChain) *Transaction {
	var inputs []TxInput
	var outputs []TxOutput

	wallets, err := wallet.CreateWallets()
	PanicHandle(err)
	w := wallets.GetWallet(from)
	publicKeyHash := wallet.PublicKeyHash(w.PublicKey)

	accumulator, validOutputs := blockchain.FindSpendableOutputs(publicKeyHash, amount)
	if accumulator < amount {
		log.Panic("ERROR: NOT ENOUGH FUNDS !")
	}
	for txID, outputs := range validOutputs {
		txIDString, err := hex.DecodeString(txID)
		PanicHandle(err)
		for _, output := range outputs {
			input := TxInput{txIDString, output, nil, w.PublicKey}
			inputs = append(inputs, input)
		}
	}
	outputs = append(outputs, *NewTxOutput(amount, to))
	if accumulator > amount {
		outputs = append(outputs, *NewTxOutput(accumulator-amount, from))
	}
	tx := Transaction{nil, inputs, outputs}
	// tx.SetID()
	tx.ID = tx.Hash()
	blockchain.SignTransaction(&tx, w.PrivateKey)
	return &tx
}

// Sign to sign the transation block to enable chaining
func (tx *Transaction) Sign(privateKey ecdsa.PrivateKey, previousTXs map[string]Transaction) {
	if tx.IsCoinBase() {
		return
	}
	for _, in := range tx.Inputs {
		if previousTXs[hex.EncodeToString(in.ID)].ID == nil {
			log.Panic("ERROR: PREVIOUS TRANSACTION DOES NOT EXIST !")
		}
	}
	txCopy := tx.TrimmedCopy()
	for inID, in := range txCopy.Inputs {
		previousTX := previousTXs[hex.EncodeToString(in.ID)]
		txCopy.Inputs[inID].Signature = nil
		txCopy.Inputs[inID].PublicKey = previousTX.Outputs[in.Out].PublicKeyHash
		txCopy.ID = txCopy.Hash()
		txCopy.Inputs[inID].PublicKey = nil
		r, s, err := ecdsa.Sign(rand.Reader, &privateKey, txCopy.ID)
		PanicHandle(err)
		signature := append(r.Bytes(), s.Bytes()...)
		tx.Inputs[inID].Signature = signature
	}
}

// Verify to verify the signature of the signed transactions
func (tx *Transaction) Verify(previousTXs map[string]Transaction) bool {
	if tx.IsCoinBase() {
		return true
	}
	for _, in := range tx.Inputs {
		if previousTXs[hex.EncodeToString(in.ID)].ID == nil {
			log.Panic("ERROR: PREVIOUS TRANSACTION DOES NOT EXIST !")
		}
	}
	txCopy := tx.TrimmedCopy()
	curve := elliptic.P256()
	for inID, in := range tx.Inputs {
		previousTX := previousTXs[hex.EncodeToString(in.ID)]
		txCopy.Inputs[inID].Signature = nil
		txCopy.Inputs[inID].PublicKey = previousTX.Outputs[in.Out].PublicKeyHash
		txCopy.ID = txCopy.Hash()
		txCopy.Inputs[inID].PublicKey = nil
		r := big.Int{}
		s := big.Int{}
		signLength := len(in.Signature)
		r.SetBytes(in.Signature[:(signLength / 2)])
		s.SetBytes(in.Signature[(signLength / 2):])
		x := big.Int{}
		y := big.Int{}
		keyLength := len(in.PublicKey)
		x.SetBytes(in.PublicKey[:(keyLength / 2)])
		y.SetBytes(in.PublicKey[(keyLength / 2):])
		rawPublicKey := ecdsa.PublicKey{Curve: curve, X: &x, Y: &y}
		if ecdsa.Verify(&rawPublicKey, txCopy.ID, &r, &s) == false {
			return false
		}
	}
	return true
}

// String to output transaction based output
func (tx Transaction) String() string {
	var lines []string
	lines = append(lines, fmt.Sprintf(" • Transaction %x:", tx.ID))
	for i, input := range tx.Inputs {
		lines = append(lines, fmt.Sprintf("   • Input %d •", i))
		lines = append(lines, fmt.Sprintf("     • Treansaction ID  : %x", input.ID))
		lines = append(lines, fmt.Sprintf("     • Out              : %d", input.Out))
		lines = append(lines, fmt.Sprintf("       • Signature      : %x", input.Signature))
		lines = append(lines, fmt.Sprintf("       • PublicKey      : %x", input.PublicKey))
	}
	for i, output := range tx.Outputs {
		lines = append(lines, fmt.Sprintf("   • Output %d •", i))
		lines = append(lines, fmt.Sprintf("     • Value            : %d", output.Value))
		lines = append(lines, fmt.Sprintf("     • Script           : %x", output.PublicKeyHash))
	}
	return strings.Join(lines, "\n")
}

// TrimmedCopy to create a trimmed copy of the entire transaction
func (tx *Transaction) TrimmedCopy() Transaction {
	var inputs []TxInput
	var outputs []TxOutput
	for _, in := range tx.Inputs {
		inputs = append(inputs, TxInput{in.ID, in.Out, nil, nil})
	}
	for _, out := range tx.Outputs {
		outputs = append(outputs, TxOutput{out.Value, out.PublicKeyHash})
	}
	txCopy := Transaction{tx.ID, inputs, outputs}
	return txCopy
}

// Hash to create a new hash for the given transaction
func (tx *Transaction) Hash() []byte {
	var hash [32]byte
	txCopy := *tx
	txCopy.ID = []byte{}
	hash = sha256.Sum256(txCopy.Serialize())
	return hash[:]
}

// IsCoinBase to check for CoinBase Transaction
func (tx *Transaction) IsCoinBase() bool {
	return len(tx.Inputs) == 1 && len(tx.Inputs[0].ID) == 0 && tx.Inputs[0].Out == -1
}

// SetID to set the ID for the Transaction
func (tx *Transaction) SetID() {
	var encoded bytes.Buffer
	var hash [32]byte
	encoder := gob.NewEncoder(&encoded)
	err := encoder.Encode(tx)
	PanicHandle(err)
	hash = sha256.Sum256(encoded.Bytes())
	tx.ID = hash[:]
}

// Serialize to serialize the transaction
func (tx *Transaction) Serialize() []byte {
	var encoded bytes.Buffer
	encoder := gob.NewEncoder(&encoded)
	err := encoder.Encode(tx)
	PanicHandle(err)
	return encoded.Bytes()
}
