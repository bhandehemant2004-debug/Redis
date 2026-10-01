package server

import (
	"fmt"
	"net"
	"os"

	"github.com/codecrafters-io/redis-starter-go/internal/client"
	"github.com/codecrafters-io/redis-starter-go/internal/protocol"
	"github.com/codecrafters-io/redis-starter-go/internal/store"
)

type  Server struct{
	port_addr string
	storage store.Storage
	role string 
	masterHost string
    masterPort string
	masterReplID       string
    masterReplOffset   int64
}
func (s* Server)Start(){
	l, err := net.Listen("tcp", s.port_addr)
	if err != nil {
	 	fmt.Println("Failed to bind to port "+s.port_addr[1:])
	 	os.Exit(1)
	}
	if s.role =="slave"{
		err := StartReplication(s)
		if err != nil {

		}
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
		masterReplID: "8371b4fb1155b71f4a04d3e1bc3e18c4a990aeeb",
        masterReplOffset: 0,
    }
}

func StartReplication(server* Server)error{
	address := net.JoinHostPort(server.masterHost,server.masterPort,)

    conn, err := net.Dial("tcp", address)
    if err != nil {
        return err
    }


    conn.Write([]byte(protocol.ResponseArray([]string{"PING"})))

	readBuffer := make([]byte, 1024)

	pending := make([]byte, 0)

	for {

		n, err := conn.Read(readBuffer)
		if err != nil {
			return err
		}
		conn.Write([]byte(protocol.ResponseArray([]string{"REPLCONF","listening-port",server.port_addr[1:]})))
		conn.Write([]byte(protocol.ResponseArray([]string{"REPLCONF","capa","psync2"})))
		pending = append(pending, readBuffer[:n]...)
		break
	}

    
	return  nil 
}