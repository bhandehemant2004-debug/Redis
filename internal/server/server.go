package server

import (
	"fmt"
	"net"
	"os"
)

type  Server struct{
	port_addr string
}

func (s* Server)Start(){
	l, err := net.Listen("tcp", s.port_addr)
	if err != nil {
	 	fmt.Println("Failed to bind to port 6379")
	 	os.Exit(1)
	}
	_, err = l.Accept()
	if err != nil {
	fmt.Println("Error accepting connection: ", err.Error())
	os.Exit(1)
}
}

func NewServer(address string) *Server {
    return &Server{
       port_addr: address,
    }
}