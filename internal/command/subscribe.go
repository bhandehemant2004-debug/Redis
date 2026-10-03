package command

import (
	"strconv"

	"github.com/codecrafters-io/redis-starter-go/internal/client"
	"github.com/codecrafters-io/redis-starter-go/internal/obj"
	"github.com/codecrafters-io/redis-starter-go/internal/store"
)


func HandleSubscribeCommand(commands []string,client *client.Client,storage store.Storage,) []byte {

	channel := commands[1]

	sub := obj.NewPubSub()

	sub.Subscribers[client] = struct{}{}

	storage.AddSubscriber(channel, client)

	return []byte(
		"*3\r\n" +
			"$9\r\nsubscribe\r\n" +
			"$" + strconv.Itoa(len(channel)) + "\r\n" +
			channel + "\r\n" +
			":1\r\n",
	)

}