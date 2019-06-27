package blockchain

import (
	"log"

	"github.com/dgraph-io/badger"
)

var (
	utxoPrefix       = []byte("utxo-")
	utxoPrefixLength = len(utxoPrefix)
)

// UTXO struct for blockchain
type UTXO struct {
	blockchain *BlockChain
}

// DeleteByPrefix to delete persistence by given Prefix
func (utx *UTXO) DeleteByPrefix(prefix []byte) {
	deleteKeys := func(keysToDelete [][]byte) error {
		if err := utx.blockchain.DataBase.Update(func(txn *badger.Txn) error {
			for _, key := range keysToDelete {
				if err := txn.Delete(key); err != nil {
					return err
				}
			}
			return nil
		}); err != nil {
			return err
		}
		return nil
	}
	collectSize := 100000
	utx.blockchain.DataBase.View(func(txn *badger.Txn) error {
		options := badger.DefaultIteratorOptions
		options.PrefetchValues = false
		iterator := txn.NewIterator(options)
		defer iterator.Close()
		keysToDelete := make([][]byte, 0, collectSize)
		keysCollected := 0
		for iterator.Seek(prefix); iterator.ValidForPrefix(prefix); iterator.Next() {
			key := iterator.Item().KeyCopy(nil)
			keysToDelete = append(keysToDelete, key)
			keysCollected++
			if keysCollected == collectSize {
				if err := deleteKeys(keysToDelete); err != nil {
					log.Panic(err)
				}
				keysToDelete = make([][]byte, 0, collectSize)
				keysCollected = 0
			}

		}
		if keysCollected > 0 {
			if err := deleteKeys(keysToDelete); err != nil {
				log.Panic(err)
			}
		}
		return nil
	})
}

// PanicHandle to Panic throw the Error
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
