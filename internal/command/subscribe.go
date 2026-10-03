package command

import (
	"fmt"

	"github.com/codecrafters-io/redis-starter-go/internal/client"
	"github.com/codecrafters-io/redis-starter-go/internal/obj"
	"github.com/codecrafters-io/redis-starter-go/internal/store"
)


func HandleSubscribeCommand(commands []string,client *client.Client,storage store.Storage,) []byte {

	channel := commands[1]

	sub := obj.NewPubSub()

	sub.Subscribers[client] = struct{}{}

	storage.AddSubscriber(channel, client)

	client.Subscriptions()[channel] = struct{}{}
	client.SetSubscribed(true)

	count := len(client.Subscriptions())

	response := fmt.Sprintf(
		"*3\r\n$9\r\nsubscribe\r\n$%d\r\n%s\r\n:%d\r\n",
		len(channel),
		channel,
		count,
	)

	return []byte(response)

}