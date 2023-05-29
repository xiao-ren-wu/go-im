package dal

import (
	"github.com/xiao-ren-wu/go-im/dal/constants"
	"net"
	"time"
)

type Server interface {
}

type Acceptor interface {
	Accept(Conn, time.Duration) (string, error)
}

type Listener interface {
	Disconnect(string) error
}

type Conn interface {
	net.Conn
	ReadFrame() (Frame, error)
	WriteFrame(constants.OpCode, []byte) error
	Flush() error
}

type MessageListener interface {
	ID() string
	Push([]byte) error
}

type Frame interface {
	SetOpCode(constants.OpCode)
	GetOpCode() constants.OpCode
	SetPayload(raw []byte)
	GetPayload() []byte
}
