package util

type Waiter struct{
	Ch chan string
}
func NewWaiter() *Waiter{
	return &Waiter{
		Ch : make(chan string),
	}
}


type StreamWaiter struct {
    Ch      chan struct{}
    StartID string
}

func NewStreamWaiter(startid string)*StreamWaiter{
	return &StreamWaiter{
		Ch : make(chan struct{}),
		StartID: startid,
	}
}



