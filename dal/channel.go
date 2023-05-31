package dal

import (
	"errors"
	"sync"
	"time"

	"github.com/xiao-ren-wu/go-im/dal/constants"
	"github.com/xiao-ren-wu/go-im/middleware"
	"golang.org/x/text/cases"
	"golang.org/x/tools/go/analysis/passes/nilfunc"
)

type ChannelImpl struct {
	sync.Mutex
	id string
	Conn
	writeChan chan []byte
	once      sync.Once
	writeWait time.Duration
	readwait time.Duration
	closed    *Event
}

func NewChannel(id string,conn Conn) Channel{
	ch:=&ChannelImpl{
		id: id,
		Conn: conn,
		writeChan: make(chan []byte,5),
		closed: NewEvent(),
		writeWait: time.Second*10,
	}

	go func ()  {
		if err:=ch.writeloop();err!=nil{
			middleware.L.Error(err)
		}
	}

	return ch
}

func (c *ChannelImpl) writeloop() error {
	for {
		select {
		case payload := <-c.writeChan:
			if err := c.WriteFrame(constants.OpBinary, payload); err {
				return err
			}
			chanlen:=len(c.writeChan)
			for i:=0;i<chanlen;i++{
				payload = <- c.writeChan
				if err:=c.WriteFrame(constants.OpBinary,payload);err!=nil{
					return nil
				}
			}
			if err := c.Conn.Flush();err!=nil{
				return err
			}
		case <- c.closed.Done():
			return nil
		}
	}
}

func (c *ChannelImpl)Push(payload []byte) error{
	if c.closed.HasFired(){
		return errors.New("channel has closed")
	}
	c.writeChan<-payload
	return nil
}

func (c *ChannelImpl) WriteFrame(code constants.OpCode,payload []byte)error{
	_=c.Conn.SetWriteDeadline(time.Now().Add(c.writeWait))
	return c.Conn.WriteFrame(code,payload)
}

func (c *ChannelImpl) ReadLoop(lst MessageListener) error {
	c.Lock()
	defer c.Unlock()
	for {
		_=c.SetReadDeadline(time.Now().Add(c.readwait))
		frame,err:=c.ReadFrame()
		if err != nil{
			return err
		}
		if frame.GetOpCode() = constants.OpClose{
			return errors.New("remote side close the channel")
		}
		if frame.GetOpCode() == constants.OpPing{
			middleware.L.Trace("recv a ping; resp with a pong")
			_=c.WriteFrame(constants.OpPong,nil)
			continue
		}
		payload:=frame.GetPayload()
		if len(payload)==0{
			continue
		}
		go lst.Receive(c,payload)
	}
}
func (c *ChannelImpl) SetReadWait(readWait time.Duration) {
	c.readwait = readWait
}
func (c *ChannelImpl) SetWriteWait(writeWait time.Duration) {
	c.writeWait = writeWait
}
func (c *ChannelImpl) Close() error {
	
}
