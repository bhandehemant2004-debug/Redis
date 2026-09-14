package command

import (
	"github.com/codecrafters-io/redis-starter-go/internal/util"

)
func HandleSetCommand(args []string,server *Server) string {
	
	obj util.Redisobj := util.Redisobj{

	}
	server.set("",)
	return "OK"
}