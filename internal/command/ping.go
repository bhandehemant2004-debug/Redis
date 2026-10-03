package command

import "github.com/codecrafters-io/redis-starter-go/internal/client"

func HandlePingCommand(c *client.Client) []byte {

	if c.IsSubscribed() {
		return []byte(
			"*2\r\n" +
				"$4\r\n" +
				"pong\r\n" +
				"$0\r\n" +
				"\r\n",
		)
	}

	return []byte("+PONG\r\n")
}