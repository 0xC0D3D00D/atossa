package main

import (
	"errors"
	"github.com/0xc0d3d00d/goresp"
	badger "github.com/dgraph-io/badger"
)

func set(args []interface{}) ([]byte, error) {
	if len(args) != 3 {
		return nil, errors.New("ERR invalid arguments")
	}

	key := args[1].([]byte)
	value := args[2].([]byte)

	err := db.Update(func(txn *badger.Txn) error {
		return txn.Set(key, value)
	})

	if err != nil {
		return nil, err
	}

	return goresp.Marshal("OK")
}

func get(args []interface{}) ([]byte, error) {
	if len(args) != 2 {
		return nil, errors.New("ERR invalid arguments")
	}

	key := args[1].([]byte)

	var value []byte
	err := db.View(func(txn *badger.Txn) error {
		item, err := txn.Get(key)
		if err != nil {
			return err
		}

		err = item.Value(func(val []byte) error {
			value = append([]byte{}, val...)
			return nil
		})
		if err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return goresp.Marshal(value)
}
