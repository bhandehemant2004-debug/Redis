package server

import (
	"net"
	"strings"

	"github.com/codecrafters-io/redis-starter-go/internal/command"
	"github.com/codecrafters-io/redis-starter-go/internal/protocol"
)


func handleconnection(conn net.Conn){
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
				result  := command.HandlePingCommand()
				response := protocol.Response_bulk_string(result)
				conn.Write([]byte(response))
			case "ECHO":
				result := command.HandleechoCommand(commands)
				response := protocol.Response_bulk_string(result)
				conn.Write([]byte(response))
			default:
				conn.Write([]byte("+PONG\r\n"))
		}
			

		}
		
	}