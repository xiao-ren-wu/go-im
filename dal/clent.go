package dal

import (
	"github.com/xiao-ren-wu/go-im/dal/constants"
	"github.com/xiao-ren-wu/go-im/dal/tcp"
	"github.com/xiao-ren-wu/go-im/middleware"
	"net"
	"time"
)

type Client interface {
	ID() string
	Name() string
	Connect(string) error
	SetDialer(Dialer)
	Send([]byte) error
	Read() (Frame, error)
	Close()
}

type Dialer interface {
	DialAndHandShake(*DialerContext) (net.Conn, error)
}

type DialerContext struct {
	Id      string
	Name    string
	Address string
	Timeout time.Duration
}

type TCPDialer struct {
	userID string
}

func (t *TCPDialer) DialAndHandShake(context *DialerContext) (net.Conn, error) {
	middleware.L.Info("begin connect %v", context.Address)
	conn, err := net.DialTimeout("TCP", context.Address, context.Timeout)
	if err != nil {
		return nil, err
	}
	if err := tcp.WriteFrame(conn, constants.OpBinary, []byte(context.Id)); err != nil {
		return nil, err
	}
	return conn, nil
}
