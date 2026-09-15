package command

import "github.com/codecrafters-io/redis-starter-go/internal/store"


func HandleLlenCommand(args []string , storage store.Storage)int{

	key := args[1];

	list := storage.Keys()

	for _,l := range list{
		if l == key{
			obj, ok := storage.Get(key)
			if !ok {
				return 0;
			}
			return  len(obj.Value.([]string))

		}
	}
	return 0
}