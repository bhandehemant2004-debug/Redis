package server

import (
	"fmt"
	"strings"
	"github.com/codecrafters-io/redis-starter-go/internal/client"
	"github.com/codecrafters-io/redis-starter-go/internal/command"
	"github.com/codecrafters-io/redis-starter-go/internal/protocol"
	"github.com/codecrafters-io/redis-starter-go/internal/rdb"
)

func handleconnection(client *client.Client, server *Server) {
	defer client.Close()

	readBuffer := make([]byte, 1024)

	pending := make([]byte, 0)

	for {

		n, err := client.Conn().Read(readBuffer)
		if err != nil {
			return
		}

		pending = append(pending, readBuffer[:n]...)

		commands, remaining, ok := protocol.Parser(pending)

		if !ok {
			continue
		}

		pending = remaining

		comm := strings.ToUpper(commands[0])
		if client.IsSubscribed() && !isAllowedInSubscribedMode(comm) {
			response := fmt.Sprintf("-ERR Can't execute '%s': only (P|S)SUBSCRIBE / (P|S)UNSUBSCRIBE / PING / QUIT / RESET are allowed in this context\r\n",strings.ToLower(comm),)

			client.WriteResponse([]byte(response))
			continue
		}

		switch comm {
		case "EXEC":
			if !client.IsMulti() {
				client.WriteResponse([]byte("-ERR EXEC without MULTI\r\n"))
				continue
			}
			if server.storage.HasWatchedKeyChanged(client.WatchedKeys()) {
				fmt.Println("WATCH conflict detected")

				client.Queue().Clear()
				clear(client.WatchedKeys())
				client.SetMulti(false)

				client.WriteResponse([]byte("*-1\r\n"))
				continue
			}
			fmt.Println("NO WATCH conflict detected")

			responses := make([][]byte, 0)

			for !client.Queue().IsEmpty() {
				queuedCommands, _ := client.Queue().Pop()
				response := Execute(queuedCommands, server)

				responses = append(responses, response)
			}
			if responses == nil {
				continue
			}

			client.SetMulti(false)
			response := protocol.ResponseTransaction(responses)

			client.WriteResponse([]byte(response))

		case "MULTI":
			if client.IsMulti() {
				client.WriteResponse([]byte("-ERR MULTI calls can not be nested\r\n"))
				continue
			}
			client.SetMulti(true)
			client.WriteResponse([]byte("+OK\r\n"))
		case "DISCARD":
			if !client.IsMulti() {
				client.WriteResponse([]byte("-ERR DISCARD without MULTI\r\n"))
				continue
			}
			for !client.Queue().IsEmpty() {
				client.Queue().Pop()
			}
			client.SetMulti(false)
			clear(client.WatchedKeys())
			client.WriteResponse([]byte("+OK\r\n"))
		case "WATCH":
			if client.IsMulti() {
				client.WriteResponse([]byte("-ERR WATCH inside MULTI is not allowed\r\n"))
				continue
			}
			for _, c := range commands[1:] {
				client.WatchedKeys()[c] = server.storage.GetVersion(c)
			}
			result := command.HandleWatchCommand(commands, server.storage)
			client.WriteResponse([]byte(result))

		case "UNWATCH":
			clear(client.WatchedKeys())
			client.WriteResponse([]byte("+OK\r\n"))

		case "REPLCONF":
			client.WriteResponse([]byte("+OK\r\n"))

		case "PSYNC":
			client.WriteResponse([]byte("+FULLRESYNC "+server.masterReplID +" 0\r\n"))
			result ,_ := rdb.EncodeRDB(&server.storage)
			rdb.SendRDB(client.Conn(),result)
		case "SUBSCRIBE":
			result := command.HandleSubscribeCommand(commands,client,server.storage)
			client.WriteResponse(result)
		default:
			if client.IsMulti() {
				client.Queue().Push(commands)
				client.WriteResponse([]byte("+QUEUED\r\n"))
				continue
			}
			response := Execute(commands, server)
			client.WriteResponse(response)
		}

	}

}

func Execute(commands []string, server *Server) []byte {
	comm := strings.ToUpper(commands[0])
	switch comm {
	
	case "INFO":
		info := fmt.Sprintf(
        "role:%s\r\nmaster_replid:%s\r\nmaster_repl_offset:%d",
        server.role,
        server.masterReplID,
        server.masterReplOffset,
    )
    	return []byte(protocol.Response_bulk_string(info))
	case "PING":
		return ([]byte("+PONG\r\n"))
	case "ECHO":
		result := command.HandleechoCommand(commands)
		response := protocol.Response_bulk_string(result)
		return ([]byte(response))
	case "SET":
		result, ok := command.HandleSetCommand(commands, server.storage)
		if ok {
			return ([]byte(result))
		}

	case "GET":
		result, ok := command.HandgetCommand(commands, server.storage)
		if !ok {
			return ([]byte("$-1\r\n"))
		}
		switch result.Type {
		case "string":
			return ([]byte(protocol.Response_bulk_string(result.Value.(string))))
		}
		return []byte("$-1\r\n")
	case "RPUSH":
		result := command.HandleRpushCommand(commands, server.storage)
		respond := protocol.AppendInteger(result)
		return ([]byte(respond))

	case "LRANGE":
		result := command.HandleLrangeCommand(commands, server.storage)
		respond := protocol.ResponseArray(result)
		return ([]byte(respond))
	case "LPUSH":
		result := command.HandleLpushCommand(commands, server.storage)
		respond := protocol.AppendInteger(result)
		return ([]byte(respond))
	case "LLEN":
		result := command.HandleLlenCommand(commands, server.storage)
		respond := protocol.AppendInteger(result)
		return ([]byte(respond))
	case "LPOP":
		result := command.HandleLpopCommand(commands, server.storage)
		if len(commands) == 2 {
			if len(result) == 0 {
				return ([]byte("$-1\r\n"))
			}
			response := protocol.Response_bulk_string(result[0])
			return ([]byte(response))

		}
		response := protocol.ResponseArray(result)
		return ([]byte(response))
	case "BLPOP":
		result := command.HandleBlpopCommand(commands, server.storage)
		if len(result) == 0 {
			return ([]byte("*-1\r\n"))
		}
		response := protocol.ResponseArray(result)
		return ([]byte(response))
	case "TYPE":
		response := command.HandleTypeCommand(commands, server.storage)
		return ([]byte(response))

	case "XADD":
		result, error := command.HandleXaddCommand(commands, server.storage)
		if error {
			response := protocol.SimpleError(result)
			return ([]byte(response))
		}

		response := protocol.Response_bulk_string(result)
		return ([]byte(response))

	case "XRANGE":
		result := command.HandleXrangeCommand(commands, server.storage)
		response := protocol.ResponseXRange(result)
		return ([]byte(response))
	case "XREAD":
		keys, collection := command.HandleXreadCommand(commands, server.storage)
		response := protocol.ResponseXRead(keys, collection)
		return ([]byte(response))

	case "INCR":
		response := command.HandleIncrCommand(commands, server.storage)
		return ([]byte(response + "\r\n"))

		
	default:
		return []byte("-ERR unknown command\r\n")
	}
	return []byte("-ERR unknown command\r\n")
}
func isAllowedInSubscribedMode(command string) bool {
	switch command {
	case "SUBSCRIBE",
		"UNSUBSCRIBE",
		"PSUBSCRIBE",
		"PUNSUBSCRIBE",
		"PING",
		"QUIT":
		return true
	default:
		return false
	}
}