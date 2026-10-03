package command

import (
	"github.com/codecrafters-io/redis-starter-go/internal/client"
	"github.com/codecrafters-io/redis-starter-go/internal/obj"
	"github.com/codecrafters-io/redis-starter-go/internal/store"
)


func HandleSubscribeCommand(commands []string,client *client.Client,storage store.Storage,) []string {

	channel := commands[1]

	sub := obj.NewPubSub()

	sub.Subscribers[client] = struct{}{}

	storage.AddSubscriber(channel, client)

	ans := []string{"subscribe", channel, "1"}


	return ans 

}