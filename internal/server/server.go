package server

import (
	"fmt"
	"net"
	"os"
	"time"

	"github.com/codecrafters-io/redis-starter-go/internal/util"
)

type  Server struct{
	port_addr string
	mpp map[string]util.Redisobj
}
func (s* Server )Set(key string ,  value util.Redisobj )(string ,bool){
	s.mpp[key]=value;
	return "+OK\r\n",true
}
func (s * Server)Get(key string )(util.Redisobj,bool){
	obj, exists := s.mpp[key]
	if !exists {
		return util.Redisobj{}, false
	}
	if obj.Expiry && obj.Expiryat.Before(time.Now()) {
			delete(s.mpp, key)
			return util.Redisobj{}, false
	}
	return obj, true
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