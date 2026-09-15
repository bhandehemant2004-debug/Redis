package command

import (
	"strconv"

	"github.com/codecrafters-io/redis-starter-go/internal/store"
)

func HandleLpopCommand(args []string, storage store.Storage)[]string{
	
	key := args[1];
	elements := 0  ;
	if len(args)>2{
		elements,_ = strconv.Atoi(args[2])
	}
	obj , ok := storage.Get(key);

	if !ok {
		return make([]string, 0)
	}
	

	list := obj.Value.([]string)

	result := make([]string,0)

	for i := 0 ;i<elements;i++{
		str := list[0];
		list = list[1:]
		result = append(result, str)
	}
	obj.Value = list

	return result
}