package constants

import "time"

type OpCode int64

const (
	OpContinuation OpCode = 0x0
	OpText         OpCode = 0x1
	OpBinary       OpCode = 0x2
	OpClose        OpCode = 0x8
	OpPing         OpCode = 0x9
	OpPong         OpCode = 0xa
)

const (
	DefaultLoginWait = 10 * time.Second
	DefauReadWait    = time.Second
	DefaultWriteWait = time.Second
)
