package command

import (
	"strconv"

	"github.com/codecrafters-io/redis-starter-go/internal/store"
)


func HandleIncrCommand(commands []string , storage store.Storage)string{



	key := commands[1]

	obj,exists := storage.Get(key)

	if exists{
		if obj.Type =="string"{
			value := obj.Value.(string)
			newvalue,done := strconv.Atoi(value)
			if done != nil{
				newvalue+=1
			}
			storage.Set(string(newvalue),obj)
			return string(newvalue)
		}

	}
	return ""
}