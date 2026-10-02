package rdb

import (
	"fmt"
	"net"

	"github.com/codecrafters-io/redis-starter-go/internal/store"
)

const (
	OpEOF      byte = 0xFF
	OpSelectDB byte = 0xFE

	TypeString byte = 0
)

func EncodeRDB(storage *store.Storage) ([]byte, error) {
	writer := NewWriter()

	writer.WriteBytes([]byte("REDIS"))
	writer.WriteBytes([]byte("0003"))


	writer.WriteByte(OpSelectDB)
	writer.WriteByte(0)


	//keys := storage.Keys()

	// for _, key := range keys {

	// 	// obj, exists := storage.Get(key)
	// 	// if !exists {
	// 	// 	continue
	// 	// }

	// 	// // value, ok := obj.(any)
	// 	// // if !ok {
	// 	// // 	return nil, fmt.Errorf(
	// 	// // 		"unsupported object type for key %q",
	// 	// // 		key,
	// 	// // 	)
	// 	// // }

	// 	// // Value type
	// 	// writer.WriteByte(TypeString)

	// 	// // Key
	// 	// writer.WriteString(key)

	// 	// // Value
	// 	// writer.WriteString(value)
	// }

	writer.WriteByte(OpEOF)

	data := writer.Bytes()




	return data, nil
}
func SendRDB(conn net.Conn, rdbData []byte) error {
	header := fmt.Sprintf("$%d\r\n", len(rdbData))

	if _, err := conn.Write([]byte(header)); err != nil {
		return err
	}

	if _, err := conn.Write(rdbData); err != nil {
		return err
	}

	return nil
}