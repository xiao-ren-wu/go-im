package websocket

import (
	"net"

	"github.com/gobwas/ws"
	"github.com/xiao-ren-wu/go-im/dal"
	"github.com/xiao-ren-wu/go-im/dal/constants"
)

type WsConn struct {
	net.Conn
}

func NewWsConn(conn net.Conn) *WsConn {
	return &WsConn{
		Conn: conn,
	}
}

func (c *WsConn) ReadFrame() (dal.Frame, error) {
	frame, err := ws.ReadFrame(c.Conn)
	if err != nil {
		return nil, err
	}
	return &Frame{raw: frame}, nil
}
func (c *WsConn) WriteFrame(opCode constants.OpCode, payload []byte) error {
	f := ws.NewFrame(ws.OpCode(opCode), true, payload)
	return ws.WriteFrame(c.Conn, f)
}
func (c *WsConn) Flush() error {
	return nil
}
