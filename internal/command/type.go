package command

import "github.com/codecrafters-io/redis-starter-go/internal/store"



func HandleTypeCommand(args [] string , storage store.Storage)string {

	key := args[1];

	obj , ok := storage.Get(key)

	if !ok {
		return "+none\r\n"
	}
	return "+"+obj.Type+"\r\n"
}