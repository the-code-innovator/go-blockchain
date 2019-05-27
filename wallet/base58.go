package wallet

import (
	"log"

	"github.com/mr-tron/base58"
)

// Base58Encode to assist in encoding the value
func Base58Encode(input []byte) []byte {
	encode := base58.Encode(input)
	return []byte(encode)
}

// Base58Decode to assist in getting the decoded value
func Base58Decode(input []byte) []byte {
	decode, err := base58.Decode(string(input[:]))
	PanicHandle(err)
	return decode

}

// PanicHandle to Panic throw errors
func PanicHandle(err error) {
	if err != nil {
		log.Panic(err)
	}
}

// ReturnHandle to return throw errors
func ReturnHandle(err error) error {
	if err != nil {
		return err
	}
	return nil
}
