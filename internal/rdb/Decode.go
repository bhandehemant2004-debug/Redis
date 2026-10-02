package rdb

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/codecrafters-io/redis-starter-go/internal/store"
)



func ReadRDBLength(reader *bufio.Reader) (int, error) {

	line, err := reader.ReadString('\n')
	if err != nil {
		return 0, err
	}

	line = strings.TrimSuffix(line, "\n")
	line = strings.TrimSuffix(line, "\r")

	if len(line) == 0 || line[0] != '$' {
		return 0, fmt.Errorf(
			"expected RDB bulk string header, got %q",
			line,
		)
	}

	length, err := strconv.Atoi(line[1:])
	if err != nil {
		return 0, fmt.Errorf(
			"invalid RDB length %q: %w",
			line[1:],
			err,
		)
	}

	if length < 0 {
		return 0, fmt.Errorf(
			"invalid RDB length: %d",
			length,
		)
	}

	return length, nil
}

func ReceiveRDB(reader *bufio.Reader) ([]byte, error) {

	length, err := ReadRDBLength(reader)
	if err != nil {
		return nil, err
	}

	rdbData := make([]byte, length)

	_, err = io.ReadFull(reader, rdbData)
	if err != nil {
		return nil, err
	}

	return rdbData, nil
}

func Decode(reader *bufio.Reader,storage *store.Storage) {
	
	_ , err :=ReceiveRDB(reader)
	if err != nil{
		return 
	}


	
	
}