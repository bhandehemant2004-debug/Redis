package command

import (
	"strconv"

	"github.com/codecrafters-io/redis-starter-go/internal/store"
)

func HandleLrangeCommand(args []string,storage store.Storage)[]string{

	key := args[1]

	obj , ok := storage.Get(key)

	if !ok {
		return make([]string, 0)
	}

	start,err := strconv.Atoi(args[2]);
	if err != nil{
		start = 0
	}
	end,err := strconv.Atoi(args[3])
	if err != nil{
		end = 0 
	}


	list := obj.Value.([]string)

	end = min(end, len(list)-1)

	if start>end{
		return make([]string, 0)
	}

	return list[start:end+1]
}