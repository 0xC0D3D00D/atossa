package main

import (
	"errors"
	"fmt"
	"github.com/0xc0d3d00d/goresp"
	badger "github.com/dgraph-io/badger/v2"
	"go.uber.org/zap"
	"io"
	"net"
	"runtime"
	"strconv"
	"strings"
)

const internalKeyPrefix = "$$$_"

var ErrWrongType = errors.New("WRONGTYPE Operation against a key holding the wrong kind of value")
var ErrIndexOutOfRange = errors.New("ERR index out of range")

type commandHandler func([]interface{}) ([]byte, error)

func UnmarshalMetadata(data []byte) (interface{}, error) {
	if len(data) == 0 {
		return nil, errors.New("Empty metadata")
	}

	metaType := data[0]
	switch metaType {
	case internalListType:
		return UnmarshalListMetadata(data)
	default:
		return nil, errors.New("Unsupported metadata type")
	}
}

func lindex(args []interface{}) ([]byte, error) {
	if len(args) != 3 {
		return nil, errors.New("ERR wrong number of arguments for 'lindex' command")
	}

	key := args[1].([]byte)
	indexStr := string(args[2].([]byte))
	index, err := strconv.ParseInt(indexStr, 10, 64)
	if err != nil {
		return nil, errors.New("ERR value is not an integer or out of range")
	}

	result, err := listIndex(key, index)
	if err != nil {
		return nil, err
	}

	if len(result) == 0 {
		return goresp.Marshal(nil)
	}

	return goresp.Marshal(result)
}

func lset(args []interface{}) ([]byte, error) {
	if len(args) != 4 {
		return nil, errors.New("ERR wrong number of arguments for 'lset' command")
	}
	key := args[1].([]byte)
	indexStr := string(args[2].([]byte))
	index, err := strconv.ParseInt(indexStr, 10, 64)
	if err != nil {
		return nil, errors.New("ERR value is not an integer or out of range")
	}
	value := args[3].([]byte)

	err = listSet(key, value, index)
	if err != nil {
		return nil, err
	}

	return goresp.Marshal("OK")
}

func lpush(args []interface{}) ([]byte, error) {
	if len(args) < 3 {
		return nil, errors.New("ERR Invalid argument")
	}

	values := [][]byte{}
	for i := 2; i < len(args); i++ {
		values = append(values, args[i].([]byte))
	}
	size, err := listPush(args[1].([]byte), values, DirectionLeft)
	if err != nil {
		return nil, err
	}
	return goresp.Marshal(size)
}

func rpush(args []interface{}) ([]byte, error) {
	if len(args) < 3 {
		return nil, errors.New("ERR Invalid argument")
	}

	values := [][]byte{}
	for i := 2; i < len(args); i++ {
		values = append(values, args[i].([]byte))
	}
	size, err := listPush(args[1].([]byte), values, DirectionRight)
	if err != nil {
		return nil, err
	}
	return goresp.Marshal(size)
}

func lpop(args []interface{}) ([]byte, error) {
	if len(args) != 2 {
		return nil, errors.New("ERR wrong number of arguments for 'lpop' command")
	}

	key := args[1].([]byte)
	value, err := listPop(key, DirectionLeft)
	if err != nil {
		return nil, err
	}
	return goresp.Marshal(value)
}

func rpop(args []interface{}) ([]byte, error) {
	if len(args) != 2 {
		return nil, errors.New("ERR wrong number of arguments for 'lpop' command")
	}

	key := args[1].([]byte)
	value, err := listPop(key, DirectionRight)
	if err != nil {
		return nil, err
	}
	return goresp.Marshal(value)
}

func llen(args []interface{}) ([]byte, error) {
	if len(args) != 2 {
		return nil, errors.New("ERR wrong number of arguments for 'llen' command")
	}

	key := args[1].([]byte)
	value, err := listLength(key)
	if err != nil {
		return nil, err
	}
	return goresp.Marshal(value)
}

func lrange(args []interface{}) ([]byte, error) {
	if len(args) != 4 {
		return nil, errors.New("ERR wrong number of arguments for 'lrange' command")
	}

	key := args[1].([]byte)

	startStr := string(args[2].([]byte))
	start, err := strconv.ParseInt(startStr, 10, 64)
	if err != nil {
		return nil, err
	}

	endStr := string(args[3].([]byte))
	end, err := strconv.ParseInt(endStr, 10, 64)
	if err != nil {
		return nil, err
	}

	values, err := listRange(key, start, end)
	if err != nil {
		return nil, err
	}

	return goresp.Marshal(values)
}

func keys(args []interface{}) ([]byte, error) {
	if len(args) != 2 {
		return nil, errors.New("ERR Invalid argument")
	}

	results := []interface{}{}
	var keyCopy []byte
	db.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.PrefetchValues = false
		it := txn.NewIterator(opts)
		defer it.Close()
		for it.Rewind(); it.Valid(); it.Next() {
			item := it.Item()
			key := item.Key()
			keyCopy = append([]byte{}, key...)
			results = append(results, keyCopy)
		}
		return nil
	})

	resp, err := goresp.Marshal(results)
	return resp, err
}

func ping(args []interface{}) ([]byte, error) {
	return goresp.Marshal("PONG")
}

func command2(args []interface{}) ([]byte, error) {
	commands := []interface{}{}
	for _, cmd := range commandMap {
		commands = append(commands, cmd.Slice())
	}
	return goresp.Marshal(commands)
}

func seq(args []interface{}) ([]byte, error) {
	result, err := db.GetSequence([]byte("_seq"), 1000)
	defer result.Release()

	id, err := result.Next()
	if err != nil {
		return nil, err
	}

	return goresp.Marshal(id)
}

func info(args []interface{}) ([]byte, error) {
	serverInfo := fmt.Sprintf("# Server\r\n"+
		"redis_version: 6.73\r\n"+
		"redis_git_sha1: 000000\r\n"+
		"redis_git_dirty: 0\r\n"+
		"redis_build_id: 1\r\n"+
		"redis_mode: standalone\r\n"+
		"os: %s\r\n"+
		"arch_bits: 64\r\n", runtime.GOOS)

	return goresp.Marshal([]byte(serverInfo))
}

func handleConnection(conn net.Conn) {
	for {
		raw, err := goresp.Unmarshal(conn)
		if err == io.EOF {
			// connection closed
			break
		}
		if err != nil {
			fmt.Printf("ERR: %s\n", err.Error())
			break
		}
		cmd := raw.([]interface{})
		if len(cmd) == 0 {
			break
		}
		cmdName := strings.ToUpper(string(cmd[0].([]byte)))
		if handler, ok := commandMap[cmdName]; ok {
			result, err := handler.handler(cmd)
			if err != nil {
				result, _ = goresp.Marshal(err.Error())
			}
			_, err = conn.Write(result)
			if err != nil {
				logger.Error("Cannot write to connection", zap.Error(err))
			}

		} else {
			logger.Info("Received unknown command", zap.String("cmd", cmdName))
			result, _ := goresp.Marshal(fmt.Errorf("ERR unknown command `%s`", cmdName))
			conn.Write(result)
		}
	}
}

var db *badger.DB
var logger *zap.Logger

func init() {
	logger, _ = zap.NewDevelopment()
	defer logger.Sync()
	var err error
	db, err = badger.Open(badger.DefaultOptions("").WithInMemory(true))
	if err != nil {
		logger.Panic(err.Error())
	}

	cmd := commandMap["COMMAND"]
	cmd.handler = command2
	commandMap["COMMAND"] = cmd
}

func main() {
	logger.Info("Artimis Server v0.1")
	listener, err := net.Listen("tcp", "0.0.0.0:6379")
	if err != nil {
		panic(err)
	}

	logger.Info("Listening on 0.0.0.0:6379")
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println(err)
		}
		go handleConnection(conn)
	}
}
