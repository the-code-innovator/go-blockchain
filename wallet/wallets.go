package wallet

import (
	"bytes"
	"crypto/elliptic"
	"encoding/gob"
	"fmt"
	"io/ioutil"
	"os"
)

const walletFile = "./tmp/wallets.data"

// Wallets structure for map of wallets
type Wallets struct {
	Wallets map[string]*Wallet
}

// CreateWallets to create a wallets file
func CreateWallets() (*Wallets, error) {
	wallets := Wallets{}
	wallets.Wallets = make(map[string]*Wallet)
	err := wallets.LoadFile()
	return &wallets, err
}

// AddWallet to add the wallet to the wallets file
func (wallets *Wallets) AddWallet() string {
	wallet := MakeWallet()
	address := fmt.Sprintf("%s", wallet.Address())
	wallets.Wallets[address] = wallet
	return address
}

// GetWallet to get the wallet
func (wallets Wallets) GetWallet(address string) Wallet {
	return *wallets.Wallets[address]
}

// GetAllAddresses to get the addresses in the wallets file
func (wallets *Wallets) GetAllAddresses() []string {
	var addresses []string
	for address := range wallets.Wallets {
		addresses = append(addresses, address)
	}
	return addresses
}

// LoadFile to load a file into the application
func (wallets *Wallets) LoadFile() error {
	if _, err := os.Stat(walletFile); os.IsNotExist(err) {
		return err
	}
	var walletsLocal Wallets
	fileContent, err := ioutil.ReadFile(walletFile)
	ReturnHandle(err)
	gob.Register(elliptic.P256())
	decoder := gob.NewDecoder(bytes.NewReader(fileContent))
	err = decoder.Decode(&walletsLocal)
	ReturnHandle(err)
	wallets.Wallets = walletsLocal.Wallets
	return nil
}

// SaveFile to save the file after edit
func (wallets *Wallets) SaveFile() {
	var content bytes.Buffer
	gob.Register(elliptic.P256())
	encoder := gob.NewEncoder(&content)
	err := encoder.Encode(wallets)
	PanicHandle(err)
	err = ioutil.WriteFile(walletFile, content.Bytes(), 0644)
	PanicHandle(err)
}
