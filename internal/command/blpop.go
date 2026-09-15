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
			obj.Value=list
			storage.Set(key,obj)
			return ans
		}
	}
	timeout := float64(0)
	if len(commands)>2{
		time, err := strconv.ParseFloat(commands[2], 64)
		if err == nil {
			timeout = time 
		}
	}
	storage.AddBlocking(key,util.NewWaiter())
	waiters,_ := storage.GetBlocking(key)	
	waiter := waiters[0]		
	if timeout ==0 {
		
		result := <-waiter.Ch
		storage.Deleteblocking(key)
		return []string{key,result}

		}

	select {
		case result := <-waiter.Ch:
			storage.Deleteblocking(key)
    		return []string{key , result}

		case <-time.After(time.Duration(timeout * float64(time.Second))):
			storage.Deleteblocking(key)
    		return []string{}
	}

	}