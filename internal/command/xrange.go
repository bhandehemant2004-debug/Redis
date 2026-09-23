package command

import (
	"github.com/codecrafters-io/redis-starter-go/internal/store"
	"github.com/codecrafters-io/redis-starter-go/internal/stream"
	"github.com/codecrafters-io/redis-starter-go/internal/obj"
)



func HandleXrangeCommand(commands []string , storage store.Storage)[]*obj.Entry{

	key := commands[1]
	start := commands[2]
	end := commands[3]

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