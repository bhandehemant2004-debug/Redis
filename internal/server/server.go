package server

import (
	"fmt"
	"net"
	"os"

	"github.com/codecrafters-io/redis-starter-go/internal/util"
)

type  Server struct{
	port_addr string
	mpp map[string]util.Redisobj
}
func (s* Server )Set(key string , value util.Redisobj ){
	s.mpp[key]=value;
}
func (s * Server)Get(key string )util.Redisobj{
	return s.mpp[key]
}
func (s* Server)Start(){
	l, err := net.Listen("tcp", s.port_addr)
	if err != nil {
	 	fmt.Println("Failed to bind to port "+s.port_addr)
	 	os.Exit(1)
	}
	
	for{
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			os.Exit(1)}
		
		go handleconnection(conn,s)
	}
	
}

func NewServer(address string) *Server {
    return &Server{
       port_addr: address,
    }
}