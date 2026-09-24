package command

import (
	"github.com/codecrafters-io/redis-starter-go/internal/obj"
	"github.com/codecrafters-io/redis-starter-go/internal/store"
	"github.com/codecrafters-io/redis-starter-go/internal/stream"
)

func HandleXreadCommand(commands []string , storage store.Storage)[]*obj.Entry{

	key := commands[2]
	start := commands[3]
	end := "+"

	if start == "-" {
		start = "0-0"
	}

	if end == "+" {
		end = "18446744073709551615-18446744073709551615"
	}
	redisobj, exists := storage.Get(key)

	if !exists{}

	tree := redisobj.Value.(*stream.Tree)

	result := make([]*obj.Entry,0)
	tree.WalkRange(nil,[]byte(start),[]byte(end),func(b []byte, a any) bool {

		entry , ok := a.(*obj.Entry)

		if !ok {
			return  false
		}
		result = append(result, entry)
		return true;
	})
	return result

}