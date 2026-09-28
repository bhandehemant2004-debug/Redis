package command

import (
	"strconv"

	"github.com/codecrafters-io/redis-starter-go/internal/store"
	"github.com/codecrafters-io/redis-starter-go/internal/util"
)

func HandleIncrCommand(commands []string, storage store.Storage) string {

	key := commands[1]

	obj, exists := storage.Get(key)

	if !exists {
		obj := util.NewRedisobj("string", "1")
		storage.Set(key, obj)
		return ":1"
	}

	if obj.Type != "string" {
		return ""
	}

	value := obj.Value.(string)

	newValue, err := strconv.Atoi(value)
	if err != nil {
		return "-ERR value is not an integer or out of range"
	}

	newValue++

	obj.Value = strconv.Itoa(newValue)

	storage.Set(key, obj)

	return ":"+strconv.Itoa(newValue)
}