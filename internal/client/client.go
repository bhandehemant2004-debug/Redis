package client

import (
	"net"

	"github.com/codecrafters-io/redis-starter-go/internal/util"
)

type Client struct {
	conn        net.Conn
	queue       util.Queue
	isMulti     bool
	watchedKeys map[string]uint64
}



func NewClient(conn net.Conn) *Client {
	return &Client{
		conn:        conn,
		isMulti:     false,
		queue:       util.NewQueue(),
		watchedKeys: make(map[string]uint64),
	}
}

func (client *Client) Close() {
	client.conn.Close()
}

func (Client *Client) WriteResponse(response []byte) {
	Client.conn.Write(response)
}

func (Client *Client) Conn() net.Conn {
	return Client.conn
}

func (c *Client) Queue() *util.Queue {
	return &c.queue
}

func (c *Client) IsMulti() bool {
	return c.isMulti
}

func (c *Client) SetMulti(value bool) {
	c.isMulti = value
}

func (c *Client) WatchedKeys() map[string]uint64 {
	return c.watchedKeys
}

// func (c *Client) HasWatchedKeyChanged(storage store.Storage) bool {
// 	for key, watchedVersion := range c.watchedKeys {
// 		currentVersion := storage.GetVersion(key)

// 		if currentVersion != watchedVersion {
// 			return true
// 		}
// 	}

// 	return false
// }
