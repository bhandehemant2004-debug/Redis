package server

import "net"

func handleconnection(conn net.Conn){
	defer conn.Close()
	buffer := make([]byte,1024)
	for{
		_,err := conn.Read(buffer)
		if err != nil{
			return
		}
		conn.Write([]byte("+PONG\r\n"))

	}
	
}