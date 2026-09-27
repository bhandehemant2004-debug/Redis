package command

import (
	"strconv"
	"strings"
	"time"

	"github.com/codecrafters-io/redis-starter-go/internal/obj"
	"github.com/codecrafters-io/redis-starter-go/internal/store"
	"github.com/codecrafters-io/redis-starter-go/internal/stream"
	"github.com/codecrafters-io/redis-starter-go/internal/util"
)

func HandleXreadCommand(commands []string, storage store.Storage) ([]string, [][]*obj.Entry) {


	command := strings.ToUpper(commands[1])
	
	
	if command == "STREAM"{
		n := (len(commands) - 2) / 2

		keys := commands[2 : 2+n]
		ids := commands[2+n:]

		allEntries := make([][]*obj.Entry, 0, len(keys))

		for i, key := range keys {

			start := ids[i]

			if start == "-" {
				start = "0-0"
			}

			end := "18446744073709551615-18446744073709551615"

			redisobj, exists := storage.Get(key)

			if !exists {
				allEntries = append(allEntries, []*obj.Entry{})
				continue
			}

			tree := redisobj.Value.(*stream.Tree)

			result := make([]*obj.Entry, 0)

			tree.WalkRange(nil,[]byte(start),[]byte(end),
				func(b []byte, a any) bool {

					entry, ok := a.(*obj.Entry)

					if !ok {
						return false
					}

					result = append(result, entry)
					return true
				},
			)

			allEntries = append(allEntries, result)
		}

		return keys, allEntries
	


	}else {

		
		timeout, err := strconv.ParseInt(commands[2], 10, 64)
		if err != nil {
			return nil, nil
		}
		key := commands[4]

		startID := commands[5]
		if startID == "-" {
			startID = "0-0"
		}

		end := "18446744073709551615-18446744073709551615"

		redisobj, exists := storage.Get(key)

		if exists {
			tree := redisobj.Value.(*stream.Tree)

			result := make([]*obj.Entry, 0)

			tree.WalkRange(nil,[]byte(startID),[]byte(end),
				func(b []byte, a any) bool {

					entry, ok := a.(*obj.Entry)

					if !ok {
						return false
					}

					result = append(result, entry)
					return true
				},
			)
			if len(result) > 0 {
				return []string{key}, [][]*obj.Entry{result}
			}
			waiter := util.NewStreamWaiter(startID)

			storage.AddStreamBlocking(key, waiter)
			timer := time.NewTimer(time.Duration(timeout) * time.Millisecond,)
			defer timer.Stop()
			select {

			case <-waiter.Ch:
				redisobj, exists := storage.Get(key)

				if !exists {
					storage.RemoveStreamBlocking(key, waiter)
					return nil, nil
				}
				tree := redisobj.Value.(*stream.Tree)

				result := make([]*obj.Entry, 0)

				tree.WalkRange(nil,[]byte(startID),[]byte(end),
					func(b []byte, a any) bool {

						entry, ok := a.(*obj.Entry)

						if !ok {
							return false
						}

						result = append(result, entry)
						return true
				},)

				storage.RemoveStreamBlocking(key, waiter)

				return []string{key}, [][]*obj.Entry{result}
			
			case <-timer.C:
				storage.RemoveStreamBlocking(key, waiter)
				return nil, nil
			}

		}
		
	}

	return nil , nil 
}