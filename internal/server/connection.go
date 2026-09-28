package server

import (
	"net"
	"strings"

	"github.com/codecrafters-io/redis-starter-go/internal/command"
	"github.com/codecrafters-io/redis-starter-go/internal/protocol"
	"github.com/codecrafters-io/redis-starter-go/internal/util"
)


func handleconnection(conn net.Conn, server *Server) {
	defer conn.Close()
	
	readBuffer := make([]byte, 1024)
	pending := make([]byte, 0)
	
	Queue := util.NewQueue()

	

	is_multi_command_active := false

	for {

		n, err := conn.Read(readBuffer)
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

		
		switch comm {
			case "EXEC":
				if !is_multi_command_active {
					conn.Write([]byte("-ERR EXEC without MULTI\r\n"))
					continue
    			}
				 responses := make([][]byte, 0)

				for !Queue.IsEmpty() {
					queuedCommands, _ := Queue.Pop()

					response := Execute(queuedCommands, server)

					responses = append(responses, response)
				}

				is_multi_command_active = false

				response := protocol.ResponseTransaction(responses)

				conn.Write([]byte(response))

			case "MULTI":
				if is_multi_command_active {
					conn.Write([]byte("-ERR MULTI calls can not be nested\r\n"))
					continue
				}
				is_multi_command_active = true
				conn.Write([]byte("+OK\r\n"))
			case "DISCARD":
				if !is_multi_command_active{
					conn.Write([]byte("ERR DISCARD without MULTI\r\n"))
				}
				for !Queue.IsEmpty() {
					 Queue.Pop()
				}
				conn.Write([]byte("+OK\r\n"))

			default:
				if is_multi_command_active {
					Queue.Push(commands)
					conn.Write([]byte("+QUEUED\r\n"))
					continue
				}
				response := Execute(commands, server)
				conn.Write(response)
		}
			

		}
		
	}

	func Execute(commands []string, server *Server)[]byte{
		comm := strings.ToUpper(commands[0])
		switch comm {
			case "PING":
				return([]byte("+PONG\r\n"))
			case "ECHO":
				result := command.HandleechoCommand(commands)
				response := protocol.Response_bulk_string(result)
				return([]byte(response))
			case "SET":
				result,ok := command.HandleSetCommand(commands,server.storage)
				if ok {
					return([]byte(result))
				}
				
			case "GET":
				result,ok := command.HandgetCommand(commands,server.storage)
				if !ok {
					return([]byte("$-1\r\n"))
				}
		
				switch result.Type{
					case "string":
					return([]byte(protocol.Response_bulk_string(result.Value.(string))))
					}
				return []byte("$-1\r\n")
			case "RPUSH":
				result := command.HandleRpushCommand(commands,server.storage)
				respond := protocol.AppendInteger(result)
				return([]byte(respond))		
				
			case "LRANGE":
				result := command.HandleLrangeCommand(commands,server.storage)
				respond := protocol.ResponseArray(result)
				return([]byte(respond))
			case "LPUSH":
				result := command.HandleLpushCommand(commands,server.storage)
				respond := protocol.AppendInteger(result)
				return([]byte(respond))
			case "LLEN":
				result := command.HandleLlenCommand(commands,server.storage)
				respond := protocol.AppendInteger(result)
				return([]byte(respond))
			case "LPOP":
				result := command.HandleLpopCommand(commands,server.storage)
				if len(commands)==2{
					if len(result)==0{
						return([]byte("$-1\r\n"))
					}
					response := protocol.Response_bulk_string(result[0])
					return([]byte(response));
					
				}
				response := protocol.ResponseArray(result)
				return([]byte(response))
			case "BLPOP":
				result := command.HandleBlpopCommand(commands,server.storage)
				if len(result)==0 {
					return([]byte("*-1\r\n"))
				}
				response := protocol.ResponseArray(result)
				return([]byte(response))
			case "TYPE":
				response := command.HandleTypeCommand(commands,server.storage)
				return([]byte(response))

			case "XADD":
				result,error := command.HandleXaddCommand(commands,server.storage)
				if error {
					response := protocol.SimpleError(result);
					return([]byte(response))
				}
				
				response := protocol.Response_bulk_string(result)
				return([]byte(response))
			
			case "XRANGE":
				result := command.HandleXrangeCommand(commands,server.storage)
				response := protocol.ResponseXRange(result)
				return([]byte(response))
			case "XREAD":
				keys,collection := command.HandleXreadCommand(commands,server.storage)
				response := protocol.ResponseXRead(keys,collection)
				return([]byte(response))

			case "INCR":
				response := command.HandleIncrCommand(commands,server.storage)
				return([]byte(response+"\r\n"))
			default:
				return([]byte("+PONG\r\n"))
		}
		return([]byte("+PONG\r\n"))
	}