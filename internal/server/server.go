package server

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"

	"github.com/codecrafters-io/redis-starter-go/internal/client"
	"github.com/codecrafters-io/redis-starter-go/internal/protocol"
	"github.com/codecrafters-io/redis-starter-go/internal/rdb"
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
		go s.StartReplication()
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

func (server *Server) StartReplication() {
    address := net.JoinHostPort(server.masterHost,server.masterPort,)

    conn, err := net.Dial("tcp", address)
    if err != nil {
        fmt.Println("Failed to connect to master:", err)
        return
    }
    defer conn.Close()

    reader := bufio.NewReader(conn)

    err = PerformHandShake(conn, reader, server)
    if err != nil {
        fmt.Println("Replication handshake failed:", err)
        return
    }

    for {
        response, err := readRESPResponse(reader)
        if err != nil {
            fmt.Println("Master connection closed:", err)
            return
        }

        fmt.Println("Received from master:", string(response))

        // Later:
        // processReplicationData(response, server)
    }
}

func PerformHandShake(conn net.Conn,reader *bufio.Reader,server *Server,) error {

    _, err := conn.Write(
        []byte(protocol.ResponseArray([]string{"PING"})),
    )
    if err != nil {
        return err
    }

    response, err := readRESPResponse(reader)
    if err != nil {
        return err
    }

    if string(response) != "+PONG\r\n" {
        return fmt.Errorf(
            "expected PONG, got %q",
            string(response),
        )
    }

    fmt.Println("Master replied PONG")


    _, err = conn.Write(
        []byte(protocol.ResponseArray([]string{
            "REPLCONF",
            "listening-port",
            server.port_addr[1:],
        })),
    )
    if err != nil {
        return err
    }

    response, err = readRESPResponse(reader)
    if err != nil {
        return err
    }

    if string(response) != "+OK\r\n" {
        return fmt.Errorf(
            "expected OK after listening-port, got %q",
            string(response),
        )
    }


    _, err = conn.Write(
        []byte(protocol.ResponseArray([]string{"REPLCONF","capa","psync2"})),)

    if err != nil {
        return err
    }

    response, err = readRESPResponse(reader)

    if err != nil {
        return err
    }

    if string(response) != "+OK\r\n" {
        return fmt.Errorf("expected OK after capa, got %q",string(response),)
    }


    _, err = conn.Write([]byte(protocol.ResponseArray([]string{"PSYNC","?","-1",})),)
    if err != nil {
        return err
    }

    response, err = readRESPResponse(reader)
    if err != nil {
        return err
    }
	err = server.HandleFullResync(response)

	if err != nil{
		return nil
	}

    fmt.Println("PSYNC response:", string(response))
    rdb.Decode(reader,&server.storage)

    return nil
}

func readRESPResponse(reader *bufio.Reader) ([]byte, error) {
    firstByte, err := reader.ReadByte()
    if err != nil {
        return nil, err
    }

    switch firstByte {
    case '+', '-':
        line, err := reader.ReadBytes('\n')
        if err != nil {
            return nil, err
        }

        return append([]byte{firstByte}, line...), nil

    default:
        return nil, fmt.Errorf(
            "unsupported RESP response type: %q",
            firstByte,
        )
    }
}
func (s *Server) HandleFullResync(response []byte) error {
    parts := strings.Fields(string(response))

    if len(parts) != 4 {
        return fmt.Errorf("invalid FULLRESYNC response: %q", string(response))
    }

    if parts[0] != "+FULLRESYNC" {
        return fmt.Errorf("expected FULLRESYNC, got: %q", string(response))
    }

    offset, err := strconv.ParseInt(parts[3], 10, 64)
    if err != nil {
        return fmt.Errorf("invalid replication offset: %q", parts[3])
    }

    s.masterReplID = parts[1]
    s.masterReplOffset = offset

    return nil
}
