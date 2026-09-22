package command

import (
	"strings"

	"github.com/codecrafters-io/redis-starter-go/internal/obj"
	"github.com/codecrafters-io/redis-starter-go/internal/store"
	"github.com/codecrafters-io/redis-starter-go/internal/stream"
	"github.com/codecrafters-io/redis-starter-go/internal/util"
)


func HandleXaddCommand(commands []string , storage store.Storage)(string,bool){
	

	key := commands[1]


	object , exists := storage.Get(key)

	if !exists{
		object = util.NewRedisobj("stream",stream.New())
	}

	tree := object.Value.(*stream.Tree)

	entry_key := commands[2]

	last_entry_key,_,ok :=tree.Last() 

	if(ok && strings.Compare(string(last_entry_key),entry_key)<=0){
		return "(error) ERR The ID specified in XADD is equal or smaller than the target stream top item",true
	}
	Entry := obj.NewEntry()

	Entry.AddEntry(commands[3:])

	tree.Insert([]byte(entry_key),Entry);


	storage.Set(key,object)



	return entry_key,false
}