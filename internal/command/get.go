package command

import (
	"github.com/codecrafters-io/redis-starter-go/internal/store"
	"github.com/codecrafters-io/redis-starter-go/internal/util"
)

func HandgetCommand(args []string, storage store.Storage) (util.Redisobj ,bool) {
	key := args[1]
	m, ok := storage.Get(key)
	if !ok {
		return util.Redisobj{}, false
	}
	return m, true
}