package main

import (
	"flag"
	"fmt"
	"net"
	"os"

	"github.com/codecrafters-io/redis-starter-go/internal/server"
)

// Ensures gofmt doesn't remove the "net" and "os" imports in stage 1 (feel free to remove this!)
var _ = net.Listen
var _ = os.Exit

func main() {
	fmt.Println("Logs from your program will appear here!")

	port := flag.String("port", "6379", "The port to bind the Redis server to")

	flag.Parse()

	address := fmt.Sprintf("0.0.0.0:%s", *port)

	s:= server.NewServer(address)

	s.Start()


}
