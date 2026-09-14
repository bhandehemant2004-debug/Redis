package server

import (
	"net"
	"strings"

	"github.com/codecrafters-io/redis-starter-go/internal/command"
	"github.com/codecrafters-io/redis-starter-go/internal/protocol"
	
)


func handleconnection(conn net.Conn,server *Server){
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
				//result  := command.HandlePingCommand()
				//response := protocol.Response_bulk_string(result)
				//conn.Write([]byte(response))
				conn.Write([]byte("+PONG\r\n"))
			case "ECHO":
				result := command.HandleechoCommand(commands)
				response := protocol.Response_bulk_string(result)
				conn.Write([]byte(response))
			case "SET":
				result , ok := server.Set(commands[1],command.HandleSetCommand(commands))
				if !ok {
					conn.Write([]byte("$-1\r\n"))
				} else {
					conn.Write([]byte(result))
				}
			case "GET":
				m,ok := server.Get(commands[1])
				if !ok {
					conn.Write([]byte("$-1\r\n"))
				} else {
					response := protocol.Response_bulk_string(m.Value)
					conn.Write([]byte(response))
				}
			default:
				conn.Write([]byte("+PONG\r\n"))
		}
			

		}
		
	}