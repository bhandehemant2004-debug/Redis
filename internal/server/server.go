package server

import (
	"fmt"
	"net"
	"os"

	"github.com/codecrafters-io/redis-starter-go/internal/client"
	"github.com/codecrafters-io/redis-starter-go/internal/store"
)

type  Server struct{
	port_addr string
	storage store.Storage
	role string 
	masterHost string
    masterPort string
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
		
		newclient := client.NewClient(conn)
		go handleconnection(newclient,s)
	}
	
}

func NewServer(address string, role string , masterHost string , masterPort string) *Server {
	storage := store.NewMemoryStorage()
    return &Server{
       port_addr: address,
	   storage: storage,
	   role: role,
	    masterHost: masterHost,
        masterPort: masterPort,
    }
}