package dal

import "time"

type ChannelMap interface {
	Add(Channel)
	Remove(id string)
	Get(id string) (Channel, bool)
	All() []Channel
}

type Agent interface {
	ID() string
	Push([]byte) error
}

type Channel interface {
	Conn
	Agent
	Close() error
	ReadLoop(MessageListener) error
	SetReadWait(time.Duration)
	SetWriteWait(time.Duration)
}
