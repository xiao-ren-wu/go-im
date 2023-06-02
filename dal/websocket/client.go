package websocket

import (
	"errors"
	"fmt"
	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
	"github.com/xiao-ren-wu/go-im/dal"
	"github.com/xiao-ren-wu/go-im/dal/constants"
	"github.com/xiao-ren-wu/go-im/middleware"
	"net"
	"net/url"
	"sync"
	"sync/atomic"
	"time"
)

type ClientOptions struct {
	Heartbeat time.Duration //登录超时
	ReadWait  time.Duration //读超时
	WriteWait time.Duration //写超时
}

type Client struct {
	sync.Mutex
	dal.Dialer
	once    sync.Once
	id      string
	name    string
	conn    net.Conn
	state   int32
	options *ClientOptions
}

func NewClient(id, name string, options *ClientOptions) dal.Client {
	if options.ReadWait == 0 {
		options.ReadWait = constants.DefaultReadWait
	}
	if options.WriteWait == 0 {
		options.WriteWait = constants.DefaultWriteWait
	}

	return &Client{
		id:      id,
		name:    name,
		options: options,
	}
}

func (c *Client) ID() string {
	return c.id
}

func (c *Client) Name() string {
	return c.name
}

func (c *Client) Connect(addr string) error {
	_, err := url.Parse(addr)
	if err != nil {
		return err
	}
	if !atomic.CompareAndSwapInt32(&c.state, 0, 1) {
		return fmt.Errorf("client has connected")
	}
	conn, err := c.Dialer.DialAndHandShake(&dal.DialerContext{
		Id:      c.id,
		Name:    c.name,
		Address: addr,
		Timeout: constants.DefaultLoginWait,
	})
	if err != nil {
		return err
	}
	if err != nil {
		atomic.CompareAndSwapInt32(&c.state, 1, 0)
		return err
	}
	if conn == nil {
		return fmt.Errorf("conn is nil")
	}
	c.conn = conn
	if c.options.Heartbeat > 0 {
		go func() {
			err := c.heartbeatLoop(conn)
			if err != nil {
				middleware.L.Error("heartbeatLoop stopped ", err)
			}
		}()
	}
	return nil
}

func (c *Client) SetDialer(dialer dal.Dialer) {
	c.Dialer = dialer
}

func (c *Client) Send(payload []byte) error {
	if atomic.LoadInt32(&c.state) == 0 {
		return fmt.Errorf("connection is nil")
	}
	c.Lock()
	defer c.Unlock()
	err := c.conn.SetWriteDeadline(time.Now().Add(c.options.WriteWait))
	if err != nil {
		return err
	}
	// 客户端消息需要使用MASK
	return wsutil.WriteClientMessage(c.conn, ws.OpBinary, payload)
}

func (c *Client) Read() (dal.Frame, error) {
	if c.conn == nil {
		return nil, errors.New("conn is nil")
	}
	if c.options.ReadWait > 0 {
		_ = c.conn.SetWriteDeadline(time.Now().Add(c.options.ReadWait))
	}
	frame, err := ws.ReadFrame(c.conn)
	if err != nil {
		return nil, err
	}
	if frame.Header.OpCode == ws.OpCode(constants.OpClose) {
		return nil, errors.New("remote set close channel")
	}
	return &Frame{raw: frame}, nil
}

func (c *Client) Close() {
	//TODO implement me
	panic("implement me")
}

func (c *Client) heartbeatLoop(conn net.Conn) error {
	tick := time.NewTicker(c.options.Heartbeat)
	for range tick.C {
		// 发送一个ping的心跳包给服务端
		if err := c.ping(conn); err != nil {
			return err
		}
	}
	return nil
}

func (c *Client) ping(conn net.Conn) error {
	c.Lock()
	defer c.Unlock()
	err := conn.SetWriteDeadline(time.Now().Add(c.options.WriteWait))
	if err != nil {
		return err
	}
	middleware.L.Trace("%s send ping to server", c.id)
	return wsutil.WriteClientMessage(conn, ws.OpPing, nil)
}
