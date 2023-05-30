package dal

import (
	"net"
	"time"
)

type Client interface {
	ID() string
	Name() string
	Connect() string
	SetDialer(Dialer)
	Send([]byte)
	Read() (Frame, error)
	Close()
}

type Dialer interface {
	DialAndHandShake(DialerContext) (net.Conn, error)
}

type DialerContext struct {
	Id      string
	Name    string
	Address string
	Timeout time.Duration
}
