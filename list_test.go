package main

import (
	badger "github.com/dgraph-io/badger/v2"
	"reflect"
	"testing"
)

func TestListPush(t *testing.T) {
	testCases := []struct {
		title        string
		key          []byte
		values       [][]byte
		direction    Direction
		result       uint32
		err          error
		dataset      [][]byte
		flushDataset bool
	}{
		{
			"nil key",
			nil,
			[][]byte{
				[]byte{'v', 'a', 'l'},
			},
			DirectionLeft,
			0,
			ErrNilKey,
			[][]byte{},
			true,
		},
		{
			"invalid direction",
			[]byte{'k', 'e', 'y'},
			[][]byte{
				[]byte{'v', 'a', 'l'},
			},
			DirectionUnknown,
			0,
			ErrInternalInvalidDirection,
			[][]byte{},
			true,
		},
		{
			"nil values",
			[]byte{'k', 'e', 'y'},
			nil,
			DirectionLeft,
			uint32(1),
			nil,
			[][]byte{
				[]byte{'$', '$', '$', '_', 'k', 'e', 'y'},
				[]byte{'L', ':', '0', ':', '0', ':', '1'},
				[]byte{'$', '$', '$', '_', 'k', 'e', 'y', ':', '0'},
				nil,
			},
			true,
		},
		{
			"single left push",
			[]byte{'k', 'e', 'y'},
			[][]byte{
				[]byte{'v', 'a', 'l'},
			},
			DirectionLeft,
			uint32(1),
			nil,
			[][]byte{
				[]byte{'$', '$', '$', '_', 'k', 'e', 'y'},
				[]byte{'L', ':', '0', ':', '0', ':', '1'},
				[]byte{'$', '$', '$', '_', 'k', 'e', 'y', ':', '0'},
				[]byte{'v', 'a', 'l'},
			},
			true,
		},
		{
			"multi left push",
			[]byte{'k', 'e', 'y'},
			[][]byte{
				[]byte{'f', 'o', 'o'},
				[]byte{'b', 'a', 'r'},
			},
			DirectionLeft,
			uint32(2),
			nil,
			[][]byte{
				[]byte{'$', '$', '$', '_', 'k', 'e', 'y'},
				[]byte{'L', ':', '0', ':', '1', ':', '2'},
				[]byte{'$', '$', '$', '_', 'k', 'e', 'y', ':', '0'},
				[]byte{'f', 'o', 'o'},
				[]byte{'$', '$', '$', '_', 'k', 'e', 'y', ':', '1'},
				[]byte{'b', 'a', 'r'},
			},
			true,
		},
		{
			"single right push",
			[]byte{'k', 'e', 'y'},
			[][]byte{
				[]byte{'v', 'a', 'l'},
			},
			DirectionRight,
			uint32(1),
			nil,
			[][]byte{
				[]byte{'$', '$', '$', '_', 'k', 'e', 'y'},
				[]byte{'L', ':', '0', ':', '0', ':', '1'},
				[]byte{'$', '$', '$', '_', 'k', 'e', 'y', ':', '0'},
				[]byte{'v', 'a', 'l'},
			},
			true,
		},
		{
			"multi right push",
			[]byte{'k', 'e', 'y'},
			[][]byte{
				[]byte{'f', 'o', 'o'},
				[]byte{'b', 'a', 'r'},
			},
			DirectionRight,
			uint32(2),
			nil,
			[][]byte{
				[]byte{'$', '$', '$', '_', 'k', 'e', 'y'},
				[]byte{'L', ':', '0', ':', '1', ':', '2'},
				[]byte{'$', '$', '$', '_', 'k', 'e', 'y', ':', '0'},
				[]byte{'f', 'o', 'o'},
				[]byte{'$', '$', '$', '_', 'k', 'e', 'y', ':', '1'},
				[]byte{'b', 'a', 'r'},
			},
			true,
		},
		{
			"consecutive left pushes #1",
			[]byte{'k', 'e', 'y'},
			[][]byte{
				[]byte{'f', 'o', 'o'},
			},
			DirectionLeft,
			uint32(1),
			nil,
			[][]byte{
				[]byte{'$', '$', '$', '_', 'k', 'e', 'y'},
				[]byte{'L', ':', '0', ':', '0', ':', '1'},
				[]byte{'$', '$', '$', '_', 'k', 'e', 'y', ':', '0'},
				[]byte{'f', 'o', 'o'},
			},
			true,
		},
		{
			"consecutive left pushes #2",
			[]byte{'k', 'e', 'y'},
			[][]byte{
				[]byte{'b', 'a', 'r'},
			},
			DirectionLeft,
			uint32(2),
			nil,
			[][]byte{
				[]byte{'$', '$', '$', '_', 'k', 'e', 'y'},
				[]byte{'L', ':', '-', '1', ':', '0', ':', '2'},
				[]byte{'$', '$', '$', '_', 'k', 'e', 'y', ':', '-', '1'},
				[]byte{'b', 'a', 'r'},
				[]byte{'$', '$', '$', '_', 'k', 'e', 'y', ':', '0'},
				[]byte{'f', 'o', 'o'},
			},
			false,
		},
		{
			"consecutive right pushes #1",
			[]byte{'k', 'e', 'y'},
			[][]byte{
				[]byte{'f', 'o', 'o'},
			},
			DirectionRight,
			uint32(1),
			nil,
			[][]byte{
				[]byte{'$', '$', '$', '_', 'k', 'e', 'y'},
				[]byte{'L', ':', '0', ':', '0', ':', '1'},
				[]byte{'$', '$', '$', '_', 'k', 'e', 'y', ':', '0'},
				[]byte{'f', 'o', 'o'},
			},
			true,
		},
		{
			"consecutive right pushes #2",
			[]byte{'k', 'e', 'y'},
			[][]byte{
				[]byte{'b', 'a', 'r'},
			},
			DirectionRight,
			uint32(2),
			nil,
			[][]byte{
				[]byte{'$', '$', '$', '_', 'k', 'e', 'y'},
				[]byte{'L', ':', '0', ':', '1', ':', '2'},
				[]byte{'$', '$', '$', '_', 'k', 'e', 'y', ':', '0'},
				[]byte{'f', 'o', 'o'},
				[]byte{'$', '$', '$', '_', 'k', 'e', 'y', ':', '1'},
				[]byte{'b', 'a', 'r'},
			},
			false,
		},
	}

	for _, testCase := range testCases {

		if testCase.flushDataset {
			err := db.DropAll()
			if err != nil {
				t.Fatal(err)
			}
		}

		actualResult, actualErr := listPush(testCase.key, testCase.values, testCase.direction)

		txn := db.NewTransaction(false)
		defer txn.Discard()
		it := txn.NewIterator(badger.DefaultIteratorOptions)
		defer it.Close()
		actualDataset := [][]byte{}
		for it.Rewind(); it.Valid(); it.Next() {
			item := it.Item()
			value, err := item.ValueCopy(nil)
			if err != nil {
				t.Fatal(err)
			}
			actualDataset = append(actualDataset, item.Key())
			actualDataset = append(actualDataset, value)
		}

		if actualResult != testCase.result || actualErr != testCase.err || !reflect.DeepEqual(actualDataset, testCase.dataset) {
			t.Fatalf("Case \"%s\":\n Expected result=%v, err=%v, dataset=%v\nActual result=%v, err=%v, dataset=%v", testCase.title, testCase.result, testCase.err, testCase.dataset, actualResult, actualErr, actualDataset)
		}
	}
}
