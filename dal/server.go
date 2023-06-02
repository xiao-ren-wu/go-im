package dal

import (
	"context"
	"github.com/segmentio/ksuid"
	"net"
	"sync"
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
	Accept(Conn, time.Duration) (string, Meta, error)
}
type Meta map[string]string

type StateListener interface {
	Disconnect(string) error
}
type DefaultAcceptor struct {
}

// Accept DefaultAcceptor
func (a *DefaultAcceptor) Accept(conn Conn, timeout time.Duration) (string, Meta, error) {
	return ksuid.New().String(), Meta{}, nil
}

type Conn interface {
	net.Conn
	ReadFrame() (Frame, error)
	WriteFrame(constants.OpCode, []byte) error
	Flush() error
}

type MessageListener interface {
	Receive(Agent, []byte)
}

type Frame interface {
	SetOpCode(constants.OpCode)
	GetOpCode() constants.OpCode
	SetPayload(raw []byte)
	GetPayload() []byte
}

type ChannelMap interface {
	Add(Channel)
	Remove(id string)
	Get(id string) (Channel, bool)
	All() []Channel
}

type ChannelMapImpl struct {
	channels *sync.Map
}

func (c *ChannelMapImpl) Add(channel Channel) {
	c.channels.Store(channel.ID(), channel)
}

func (c *ChannelMapImpl) Remove(id string) {
	c.channels.Delete(id)
}

func (c *ChannelMapImpl) Get(id string) (Channel, bool) {
	value, ok := c.channels.Load(id)
	if ok {
		return value.(Channel), ok
	}
	return nil, false
}

func (c *ChannelMapImpl) All() (resList []Channel) {
	c.channels.Range(func(key, value any) bool {
		resList = append(resList, value.(Channel))
		return true
	})
	return
}

func NewChannelMap(size int) ChannelMap {
	return &ChannelMapImpl{channels: new(sync.Map)}
}

type Agent interface {
	ID() string
	Push([]byte) error
}

type Channel interface {
	Conn
	Agent
	Close() error
	Readloop(MessageListener) error
	SetReadWait(time.Duration)
	SetWriteWait(time.Duration)
}
