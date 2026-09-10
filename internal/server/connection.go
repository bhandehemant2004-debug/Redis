package server

import "net"

func handleconnection(conn net.Conn){
	conn.Write([]byte("+PONG\r\n"))
}