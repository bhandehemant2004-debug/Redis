package command

import (
	"strconv"

	"github.com/codecrafters-io/redis-starter-go/internal/client"
	"github.com/codecrafters-io/redis-starter-go/internal/store"
)

func HandlePublisCommand(commands []string , client *client.Client, storage store.Storage)[]byte{
	channel := commands[1]
	message := commands[2]
	sub, exists := storage.GetSubscribers(channel)

	if !exists {
		return []byte(":0\r\n")
	}
	response := []byte(
		"*3\r\n" +
			"$7\r\n" +
			"message\r\n" +
			"$" + strconv.Itoa(len(channel)) + "\r\n" +
			channel + "\r\n" +
			"$" + strconv.Itoa(len(message)) + "\r\n" +
			message + "\r\n",
	)

	for subscriber := range sub.Subscribers {
		subscriber.WriteResponse(response)
	}

	count := len(sub.Subscribers)

	return []byte(":" + strconv.Itoa(count) + "\r\n")
}