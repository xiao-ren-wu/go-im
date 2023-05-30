package dal

import (
	"context"
	"net"
	"time"

	"github.com/xiao-ren-wu/go-im/dal/constants"
)

type Server interface {
	SetAcceptor(Acceptor)
	SetMessageListener(MessageListener)
	SetStateListener(StateListener)
	SetReadWait(time.Duration)
	SetChannelMap(ChannelMap)
	Start() error
	Push(string, []byte) error
	Shutdown(context.Context) error
}

type Acceptor interface {
	Accept(Conn, time.Duration) (string, error)
}

type StateListener interface {
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
