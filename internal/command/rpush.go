package command

import (
	"github.com/codecrafters-io/redis-starter-go/internal/store"
	"github.com/codecrafters-io/redis-starter-go/internal/util"
)


func HandleRpushCommand(args []string, storage store.Storage)int {
    key := args[1]

    obj, ok := storage.Get(key)

    if !ok {
        obj = util.Redisobj{
            Type:  "list",
            Value: []string{},
        }
    }

    list := obj.Value.([]string)

    for _, val := range args[2:] {
        list = append(list, val)
    }

    obj.Value = list

    storage.Set(key, obj)
    waiters,status := storage.GetBlocking(key)

    if status{
        waiter := waiters[0]
        waiter.Ch<-obj.Value.([]string)[0]
    }

	return len(list)
}