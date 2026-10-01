package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"strings"

	"github.com/codecrafters-io/redis-starter-go/internal/server"
)

// Ensures gofmt doesn't remove the "net" and "os" imports in stage 1 (feel free to remove this!)
var _ = net.Listen
var _ = os.Exit

func main() {
	fmt.Println("Logs from your program will appear here!")

	port := flag.String("port", "6379", "server port")

    replicaOf := flag.String("replicaof", "", "master host and port")

    flag.Parse()

    role := "master"
    masterHost := ""
    masterPort := ""

    if *replicaOf != "" {
        parts := strings.Fields(*replicaOf)

        if len(parts) != 2 {
            log.Fatal("invalid --replicaof: expected <host> <port>")
        }

        role = "slave"
        masterHost = parts[0]
        masterPort = parts[1]
    }

	s:= server.NewServer(":"+*port,role,masterHost,masterPort)

	s.Start()


}
