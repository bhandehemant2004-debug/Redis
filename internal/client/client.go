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
	subscriptions map[string]struct{}
	subscribed bool
}



func NewClient(conn net.Conn) *Client {
	return &Client{
		conn:        conn,
		isMulti:     false,
		queue:       util.NewQueue(),
		watchedKeys: make(map[string]uint64),
		subscriptions: make(map[string]struct{}),
		subscribed: false,
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

func (c *Client) Subscriptions() map[string]struct{} {
	return c.subscriptions
}

func (c *Client)IsSubscribed()bool{
	return c.subscribed
}

func (c*Client)SetSubscribed(is bool){
	c.subscribed=is
}
