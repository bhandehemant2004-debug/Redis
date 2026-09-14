package command

import (
	"strconv"
	"time"

	"github.com/codecrafters-io/redis-starter-go/internal/util"
)

func HandleSetCommand(args []string) (util.Redisobj) {
	value := args[2]
	if len(args) > 3 && args[3] == "EX" {
		seconds , err := strconv.Atoi(args[4])
		if err != nil {
			return util.Redisobj{
				Value:   value,
				Expiryat: time.Now().Add(time.Duration(seconds) * time.Second),
				Expiry:  true,
			}
		}
	}	
	if len(args) > 3 && args[3] == "PX" {
		milliseconds , err := strconv.Atoi(args[4])
		if err != nil {
			return util.Redisobj{
				Value:   value,
				Expiryat: time.Now().Add(time.Duration(milliseconds) * time.Millisecond),
				Expiry:  true,
			}
		}
	}
	obj := util.Redisobj{
		Value: value,
	}
	return obj	
}
