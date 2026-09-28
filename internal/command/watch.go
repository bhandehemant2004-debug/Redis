package command

import "github.com/codecrafters-io/redis-starter-go/internal/store"


func HandleWatchCommand(commands []string , storage store.Storage)string{
	return "+OK\r\n"
}