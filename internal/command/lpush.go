package command

import (
	"github.com/codecrafters-io/redis-starter-go/internal/store"
	"github.com/codecrafters-io/redis-starter-go/internal/util"
)


func HandleLpushCommand(args []string,storage store.Storage) int {
	key := args[1]

    obj, ok := storage.Get(key)

    if !ok {
        obj = util.Redisobj{
            Type:  "list",
            Value: []string{},
        }
    }

    list := obj.Value.([]string)
    newlist := make([]string,0)

	leng := len(args)-1;

	for i := leng ; i>=2 ;i--{
		newlist = append(newlist, args[i])
	}
    for _, val := range list {
        newlist = append(newlist, val)
    }

    obj.Value = newlist

    storage.Set(key, obj)

	return len(newlist)

}