package command

import (
	"github.com/codecrafters-io/redis-starter-go/internal/obj"
	"github.com/codecrafters-io/redis-starter-go/internal/store"
	"github.com/codecrafters-io/redis-starter-go/internal/stream"
)

func HandleXreadCommand(commands []string, storage store.Storage) ([]string, [][]*obj.Entry) {

	// commands:
	//
	// XREAD STREAMS key1 key2 id1 id2
	//
	// Find number of streams.
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

		tree.WalkRange(
			nil,
			[]byte(start),
			[]byte(end),
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
}