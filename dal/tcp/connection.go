package tcp

import (
	"github.com/xiao-ren-wu/go-im/dal/constants"
	"github.com/xiao-ren-wu/go-im/dal/wire/endian"
	"io"
)

func WriteFrame(w io.Writer, code constants.OpCode, payload []byte) error {
	if err := endian.WriteUint8(w, uint8(code)); err != nil {
		return err
	}
	if err := endian.WriteBytes(w, payload); err != nil {
		return err
	}
	return nil
}
