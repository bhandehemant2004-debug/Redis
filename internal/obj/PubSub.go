package obj

import "github.com/codecrafters-io/redis-starter-go/internal/client"

type PubSub struct {
    Subscribers map[*client.Client]struct{}
}

func NewPubSub()*PubSub{
	return &PubSub{
		Subscribers: make(map[*client.Client]struct{}),
	}
}