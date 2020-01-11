package main

import (
	"errors"
	"fmt"
	badger "github.com/dgraph-io/badger"
	"strconv"
	"strings"
)

const internalListType = 'L'

var ErrInvalidListMetadata = errors.New("Invalid list metadata")

type Direction uint8

const (
	DirectionUnknown Direction = iota
	DirectionLeft
	DirectionRight
)

type ListMetadata struct {
	first int64
	last  int64
	size  uint32
}

func (lm ListMetadata) String() string {
	firstStr := strconv.FormatInt(lm.first, 16)
	lastStr := strconv.FormatInt(lm.last, 16)
	size := strconv.FormatUint(uint64(lm.size), 16)

	str := append([]byte{internalListType}, ':')
	str = append(str, firstStr...)
	str = append(str, ':')
	str = append(str, lastStr...)
	str = append(str, ':')
	str = append(str, size...)

	return string(str)
}

func UnmarshalListMetadata(data []byte) (interface{}, error) {
	dataStr := string(data)
	params := strings.Split(dataStr, ":")
	if len(params) != 4 {
		return nil, ErrInvalidListMetadata
	}
	first, err := strconv.ParseInt(params[1], 16, 64)
	if err != nil {
		return nil, ErrInvalidListMetadata
	}
	last, err := strconv.ParseInt(params[2], 16, 64)
	if err != nil {
		return nil, ErrInvalidListMetadata
	}
	size, err := strconv.ParseUint(params[3], 16, 64)
	if err != nil {
		return nil, ErrInvalidListMetadata
	}

	return ListMetadata{first, last, uint32(size)}, nil
}

func listPop(key []byte, direction Direction) ([]byte, error) {
	internalKey := append([]byte(internalKeyPrefix), key...)

	var value []byte = nil
	err := db.Update(func(txn *badger.Txn) error {
		// Ensure there is no simple string key with the same name exists
		_, err := txn.Get(key)
		if err != badger.ErrKeyNotFound {
			return ErrWrongType
		}

		metadataItem, err := txn.Get(internalKey)
		if err == badger.ErrKeyNotFound {
			return nil
		} else if err != nil {
			return err
		}
		metadataRaw, err := metadataItem.ValueCopy(nil)
		if err != nil {
			return err
		}
		metadata, err := UnmarshalMetadata(metadataRaw)
		if err != nil {
			return err
		}
		listMetadata, ok := metadata.(ListMetadata)
		if !ok {
			return ErrWrongType
		}

		var itemId string
		if direction == DirectionLeft {
			itemId = strconv.FormatInt(listMetadata.first, 16)
			listMetadata.first++
			listMetadata.size--
		} else if direction == DirectionRight {
			itemId = strconv.FormatInt(listMetadata.last, 16)
			listMetadata.last--
			listMetadata.size--
		} else {
			return errors.New("BUG: Invalid direction")
		}

		itemKey := append([]byte{}, internalKey...)
		itemKey = append(itemKey, ':')
		itemKey = append(itemKey, itemId...)
		item, err := txn.Get(itemKey)
		if err != nil {
			return err
		}

		value, err = item.ValueCopy(nil)
		if err != nil {
			return err
		}
		err = txn.Delete(itemKey)
		if err != nil {
			return err
		}

		if listMetadata.size == 0 {
			err = txn.Delete(internalKey)
			return err
		}

		err = txn.Set(internalKey, []byte(listMetadata.String()))
		return err
	})

	return value, err
}

func listPush(key []byte, values [][]byte, direction Direction) (uint32, error) {
	internalKey := append([]byte(internalKeyPrefix), key...)

	size := uint32(0)

	err := db.Update(func(txn *badger.Txn) error {
		// Ensure there is no simple string key with the same name exists
		_, err := txn.Get(key)
		if err != badger.ErrKeyNotFound {
			return ErrWrongType
		}

		size = 1
		metadata := ListMetadata{0, 0, size}

		item, err := txn.Get(internalKey)
		if err == badger.ErrKeyNotFound {
			var value []byte
			if direction == DirectionLeft {
				value = values[len(values)-1]
				values = values[:len(values)-1]
			} else if direction == DirectionRight {
				value = values[0]
				values = values[1:]
			} else {
				return errors.New("BUG: invalid direction")
			}

			err = txn.Set(internalKey, []byte(metadata.String()))
			if err != nil {
				return err
			}

			itemKey := append([]byte{}, internalKey...)
			itemKey = append(itemKey, ":0"...)
			err = txn.Set(itemKey, value)
			if err != nil {
				return err
			}
			if len(values) == 1 {
				return nil
			}
		} else if err != nil {
			return err
		} else {
			metadataVal, err := item.ValueCopy(nil)
			if err != nil {
				return err
			}
			readMetadata, err := UnmarshalMetadata(metadataVal)
			if err != nil {
				return err
			}
			var ok bool
			metadata, ok = readMetadata.(ListMetadata)
			if !ok {
				return ErrWrongType
			}
		}

		var start, step int
		var condition func(int) bool
		if direction == DirectionLeft {
			start = len(values) - 1
			step = -1
			condition = func(i int) bool { return i >= 0 }
		} else if direction == DirectionRight {
			start = 0
			step = 1
			condition = func(i int) bool { return i < len(values) }
		} else {
			return errors.New("BUG: Invalid direction")
		}
		for i := start; condition(i); i += step {
			var itemId string
			if direction == DirectionLeft {
				metadata.first--
				itemId = strconv.FormatInt(metadata.first, 16)
			} else if direction == DirectionRight {
				metadata.last++
				itemId = strconv.FormatInt(metadata.last, 16)
			} else {
				return errors.New("BUG: Invalid direction")
			}
			metadata.size++

			size = metadata.size

			itemKey := append([]byte{}, internalKey...)
			itemKey = append(itemKey, ':')
			itemKey = append(itemKey, itemId...)

			err = txn.Set(itemKey, values[i])
			if err != nil {
				return err
			}
		}
		err = txn.Set(internalKey, []byte(metadata.String()))

		return err
	})

	return size, err
}

func listIndex(key []byte, index int64) ([]byte, error) {
	internalKey := append([]byte(internalKeyPrefix), key...)

	var val []byte

	fmt.Println(val)
	err := db.View(func(txn *badger.Txn) error {
		// Ensure there is no simple string key with the same name exists
		_, err := txn.Get(key)
		if err != badger.ErrKeyNotFound {
			return ErrWrongType
		}

		item, err := txn.Get(internalKey)
		if err == badger.ErrKeyNotFound {
			return errors.New("ERR no such key")
		} else if err != nil {
			return err
		}
		metadataVal, err := item.ValueCopy(nil)
		if err != nil {
			return err
		}
		readMetadata, err := UnmarshalMetadata(metadataVal)
		if err != nil {
			return err
		}
		var ok bool
		metadata, ok := readMetadata.(ListMetadata)
		if !ok {
			return ErrWrongType
		}

		if index >= int64(metadata.size) || index < -int64(metadata.size) {
			// index out of range
			return nil
		}

		if index < 0 {
			index = metadata.last + index + 1
		} else {
			index = metadata.first + index
		}

		itemId := strconv.FormatInt(index, 16)

		itemKey := append([]byte{}, internalKey...)
		itemKey = append(itemKey, ':')
		itemKey = append(itemKey, itemId...)
		item, err = txn.Get(itemKey)
		if err != nil {
			return err
		}

		val, err = item.ValueCopy(nil)
		return err
	})

	if err != nil {
		return nil, err
	}

	return val, nil
}

func listSet(key, value []byte, index int64) error {
	internalKey := append([]byte(internalKeyPrefix), key...)

	err := db.Update(func(txn *badger.Txn) error {
		// Ensure there is no simple string key with the same name exists
		_, err := txn.Get(key)
		if err != badger.ErrKeyNotFound {
			return ErrWrongType
		}

		item, err := txn.Get(internalKey)
		if err == badger.ErrKeyNotFound {
			return errors.New("ERR no such key")
		} else if err != nil {
			return err
		}
		metadataVal, err := item.ValueCopy(nil)
		if err != nil {
			return err
		}
		readMetadata, err := UnmarshalMetadata(metadataVal)
		if err != nil {
			return err
		}
		var ok bool
		metadata, ok := readMetadata.(ListMetadata)
		if !ok {
			return ErrWrongType
		}

		if index >= int64(metadata.size) || index < -int64(metadata.size) {
			return ErrIndexOutOfRange
		}

		if index < 0 {
			index = metadata.last + index + 1
		} else {
			index = metadata.first + index
		}

		itemId := strconv.FormatInt(index, 16)

		itemKey := append([]byte{}, internalKey...)
		itemKey = append(itemKey, ':')
		itemKey = append(itemKey, itemId...)
		err = txn.Set(itemKey, value)

		return err
	})

	return err
}
