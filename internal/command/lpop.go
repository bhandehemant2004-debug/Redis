package command

import (
	"strconv"

	"github.com/codecrafters-io/redis-starter-go/internal/store"
)

func HandleLpopCommand(args []string, storage store.Storage)[]string{
	
	key := args[1];
	elements := 1  ;
	if len(args)>2{
		elements,_ = strconv.Atoi(args[2])
	}
	obj , ok := storage.Get(key);

	if !ok {
		return make([]string, 0)
	}
	

	list := obj.Value.([]string)

	result := make([]string,0)
	elements = min(elements,len(list))
	for i := 0 ;i<elements;i++{
		result = append(result,list[i])
	}
	list = list[elements:]
	obj.Value = list
	storage.Set(key,obj)
	return result
}