package command

import (
	"strconv"
	"time"
	"github.com/codecrafters-io/redis-starter-go/internal/store"
	"github.com/codecrafters-io/redis-starter-go/internal/util"
)

func HandleBlpopCommand(commands []string , storage store.Storage)[]string{

	key := commands[1];

	obj , exist := storage.Get(key)
	ans := make([]string,0)
	if exist{
		list  := obj.Value.([]string)
		if len(list)>0{
			ans = append(ans, list[0])
			list = list[1:]
			obj.Value=ans
			storage.Set(key,obj)
			return ans
		}
	}
	timeout := 0 
	if len(commands)>2{
		timeout,_ = strconv.Atoi(commands[2])
	}
	storage.AddBlocking(key,util.NewWaiter())
	waiters,_ := storage.GetBlocking(key)	
	waiter := waiters[0]		
	if timeout ==0 {
		
		result := <-waiter.Ch
		storage.Deleteblocking(key)
		return []string{result}

		}

	select {
		case result := <-waiter.Ch:
			storage.Deleteblocking(key)
    		return []string{result}

		case <-time.After(time.Duration(timeout) * time.Second):
			storage.Deleteblocking(key)
    		return []string{}
	}

	}