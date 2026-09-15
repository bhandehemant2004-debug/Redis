package command

import (
	"strconv"
	"time"

	"github.com/codecrafters-io/redis-starter-go/internal/store"
	"github.com/codecrafters-io/redis-starter-go/internal/util"
)

func HandleSetCommand(args []string,storage store.Storage)( string , bool){

	key := args[1]
	value := args[2]
	if len(args) > 4 && args[3] == "EX" {
		seconds , err := strconv.Atoi(args[4])
		if err == nil {
			obj := util.Redisobj{
				Type: "string",
				Value:   value,
				Expiryat: time.Now().Add(time.Duration(seconds) * time.Second),
				Expiry:  true,
			}
			storage.Set(key, obj)
			return "+OK\r\n", true	
		}	
	}	
	if len(args) > 4 && args[3] == "PX" {
		milliseconds , err := strconv.Atoi(args[4])
		if err == nil {
			obj := util.Redisobj{
				Type: "string",
				Value:   value,
				Expiryat: time.Now().Add(time.Duration(milliseconds) * time.Millisecond),
				Expiry:  true,
			}
			storage.Set(key, obj)
			return "+OK\r\n", true
		}
	}
	
	obj := util.Redisobj{
		Value: value,
		Type: "string",
	}
	storage.Set(key, obj)
	return "+OK\r\n", true	
	
}
