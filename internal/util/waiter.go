package util

type Waiter struct{
	Ch chan string
}
func NewWaiter() *Waiter{
	return &Waiter{
		Ch : make(chan string),
	}
}