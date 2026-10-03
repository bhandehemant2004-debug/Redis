package command

import (
	"strconv"

	"github.com/codecrafters-io/redis-starter-go/internal/client"
	"github.com/codecrafters-io/redis-starter-go/internal/store"
)


func HandleUnsubscribeCommand(commands []string,client *client.Client,storage store.Storage,) []byte {

	var response []byte

	for _, channel := range commands[1:] {

		
		storage.RemoveSubscriber(channel, client)

		
		delete(client.Subscriptions(), channel)


		count := len(client.Subscriptions())

		
		response = append(response,[]byte("*3\r\n")...,)

		response = append(response,[]byte("$11\r\nunsubscribe\r\n")...,)

		response = append(response,[]byte("$"+strconv.Itoa(len(channel))+"\r\n")...,)

		response = append(response,[]byte(channel+"\r\n")...,)

		response = append(response,[]byte(":"+strconv.Itoa(count)+"\r\n")...,)
	}
	if len(client.Subscriptions())==0{
		client.SetSubscribed(false)
	}
	
	return response
}