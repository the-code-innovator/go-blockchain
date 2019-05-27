package wallet

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"

	"golang.org/x/crypto/ripemd160"
)

// constants for version handling the blockchain
const (
	checkSumLength = 4
	version        = byte(0x00)
)

// Wallet structure for the Wallet type in the blockchain
type Wallet struct {
	PrivateKey ecdsa.PrivateKey
	PublicKey  []byte
}

// NewKeyPair to create a new KeyPair
func NewKeyPair() (ecdsa.PrivateKey, []byte) {
	curve := elliptic.P256()
	private, err := ecdsa.GenerateKey(curve, rand.Reader)
	PanicHandle(err)
	public := append(private.PublicKey.X.Bytes(), private.PublicKey.Y.Bytes()...)
	return *private, public
}

// PublicKeyHash to create the public hash
func PublicKeyHash(publicKey []byte) []byte {
	publicKeyHash := sha256.Sum256(publicKey)
	ripemd160Hasher := ripemd160.New()
	_, err := ripemd160Hasher.Write(publicKeyHash[:])
	PanicHandle(err)
	publicRipeMD160Hash := ripemd160Hasher.Sum(nil)
	return publicRipeMD160Hash
}

// GenerateCheckSum to generate the checksum
func GenerateCheckSum(payload []byte) []byte {
	firstHash := sha256.Sum256(payload)
	secondHash := sha256.Sum256(firstHash[:])
	return secondHash[:checkSumLength]
}

// MakeWallet to create a wallet
func MakeWallet() *Wallet {
	private, public := NewKeyPair()
	wallet := Wallet{private, public}
	return &wallet
}

// Address to find the address of the Wallet
func (w Wallet) Address() []byte {
	publicKeyHash := PublicKeyHash(w.PublicKey)
	versionedHash := append([]byte{version}, publicKeyHash...)
	checkSum := GenerateCheckSum(versionedHash)
	fullHash := append(versionedHash, checkSum...)
	address := Base58Encode(fullHash)
	return address
}

// ValidateAddress to validate the address that is passed into the blockchain
func ValidateAddress(address string) bool {
	publicKeyHash := Base58Decode([]byte(address))
	actualChecksum := publicKeyHash[len(publicKeyHash)-checkSumLength:]
	actualVersion := publicKeyHash[0]
	publicKeyHash = publicKeyHash[1 : len(publicKeyHash)-checkSumLength]
	targetChecksum := GenerateCheckSum(append([]byte{actualVersion}, publicKeyHash...))
	return bytes.Compare(actualChecksum, targetChecksum) == 0
}
