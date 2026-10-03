package command

import (
	"github.com/codecrafters-io/redis-starter-go/internal/client"
	"github.com/codecrafters-io/redis-starter-go/internal/store"
)

func HandlePublisCommand(commands []string , client *client.Client, storage store.Storage)byte[]{
	channel := commands[1]

	sub, exists := storage.GetSubscribers(channel)

	if !exists {
		return []byte(":0\r\n")
	}

	count := len(sub.Subscribers)

	return []byte(":" + strconv.Itoa(count) + "\r\n")
}