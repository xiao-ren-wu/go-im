package websocket

import (
	"github.com/gobwas/ws"
	"github.com/xiao-ren-wu/go-im/dal/constants"
)

type Frame struct {
	raw ws.Frame
}

func (f *Frame) SetOpCode(opCode constants.OpCode) {
	f.raw.Header.OpCode = ws.OpCode(opCode)
}

func (f *Frame) GetOpCode() constants.OpCode {
	return constants.OpCode(f.raw.Header.OpCode)
}

func (f *Frame) SetPayload(raw []byte) {
	f.raw.Payload = raw
}

func (f *Frame) GetPayload() []byte {
	if f.raw.Header.Masked {
		ws.Cipher(f.raw.Payload, f.raw.Header.Mask, 0)
	}
	f.raw.Header.Masked = false
	return f.raw.Payload
}
