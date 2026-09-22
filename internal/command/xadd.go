package command

import (
	"github.com/codecrafters-io/redis-starter-go/internal/obj"
	"github.com/codecrafters-io/redis-starter-go/internal/store"
	"github.com/codecrafters-io/redis-starter-go/internal/stream"
	"github.com/codecrafters-io/redis-starter-go/internal/util"
)


func HandleXaddCommand(commands []string , storage store.Storage)string{
	

	key := commands[1]


	object , exists := storage.Get(key)

	if !exists{
		object = util.NewRedisobj("stream",stream.New())
	}

	tree := object.Value.(*stream.Tree)

	entry_key := commands[2]

	Entry := obj.NewEntry()

	Entry.AddEntry(commands[3:])

	tree.Insert([]byte(entry_key),Entry);


	storage.Set(key,object)



	return key
}