package server

import (
	"fmt"
	"net"
	"strings"

	"github.com/codecrafters-io/redis-starter-go/internal/command"
	"github.com/codecrafters-io/redis-starter-go/internal/protocol"
)


func handleconnection(conn net.Conn, server *Server) {
	defer conn.Close()
	
	readBuffer := make([]byte, 1024)
	pending := make([]byte, 0)
	

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
			case "PING":
				conn.Write([]byte("+PONG\r\n"))
			case "ECHO":
				result := command.HandleechoCommand(commands)
				response := protocol.Response_bulk_string(result)
				conn.Write([]byte(response))
			case "SET":
				result,ok := command.HandleSetCommand(commands,server.storage)
				if ok {
					conn.Write([]byte(result))
				}
				
			case "GET":
				result,ok := command.HandgetCommand(commands,server.storage)
				if !ok {
					conn.Write([]byte("$-1\r\n"))
				}
		
				switch result.Type{
					case "string":
					conn.Write([]byte(protocol.Response_bulk_string(result.Value.(string))))
					}
			
			case "RPUSH":
				result := command.HandleRpushCommand(commands,server.storage)
				respond := protocol.AppendInteger(result)
				conn.Write([]byte(respond))		
				
			case "LRANGE":
				result := command.HandleLrangeCommand(commands,server.storage)
				respond := protocol.ResponseArray(result)
				conn.Write([]byte(respond))
			case "LPUSH":
				result := command.HandleLpushCommand(commands,server.storage)
				respond := protocol.AppendInteger(result)
				conn.Write([]byte(respond))
			case "LLEN":
				result := command.HandleLlenCommand(commands,server.storage)
				respond := protocol.AppendInteger(result)
				conn.Write([]byte(respond))
			case "LPOP":
				result := command.HandleLpopCommand(commands,server.storage)
				if len(commands)==2{
					if len(result)==0{
						conn.Write([]byte("$-1\r\n"))
						continue
					}
					response := protocol.Response_bulk_string(result[0])
					conn.Write([]byte(response));
					continue
					
				}
				response := protocol.ResponseArray(result)
				conn.Write([]byte(response))
			case "BLPOP":
				result := command.HandleBlpopCommand(commands,server.storage)
				if len(result)==0 {
					conn.Write([]byte("*-1\r\n"))
					continue
				}
				response := protocol.ResponseArray(result)
				conn.Write([]byte(response))
			case "TYPE":
				response := command.HandleTypeCommand(commands,server.storage)
				conn.Write([]byte(response))

			case "XADD":
				result,error := command.HandleXaddCommand(commands,server.storage)
				if error {
					response := protocol.SimpleError(result);
					fmt.Println(response)
					conn.Write([]byte(response))
					continue
				}
				
				response := protocol.Response_bulk_string(result)
				conn.Write([]byte(response))
			
			case "XRANGE":
				result := command.HandleXrangeCommand(commands,server.storage)
				response := protocol.ResponseXRange(result)
				conn.Write([]byte(response))
			case "XREAD":
				keys,collection := command.HandleXreadCommand(commands,server.storage)
				response := protocol.ResponseXRead(keys,collection)
				conn.Write([]byte(response))

			case "INCR":
				response := command.HandleIncrCommand(commands,server.storage)
				conn.Write([]byte(":"+response+"\r\n"))
			default:
				conn.Write([]byte("+PONG\r\n"))
		}
			

		}
		
	}