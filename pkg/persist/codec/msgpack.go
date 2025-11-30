package codec

import (
	"github.com/vmihailenco/msgpack/v5"
)

// MsgPack implements Codec using MessagePack serialization.
// MessagePack is a fast, compact binary format ideal for game state.
type MsgPack struct{}

// NewMsgPack creates a new MessagePack codec.
func NewMsgPack() *MsgPack {
	return &MsgPack{}
}

// Encode serializes a value to MessagePack bytes.
func (c *MsgPack) Encode(v any) ([]byte, error) {
	if v == nil {
		return nil, ErrNilValue
	}
	return msgpack.Marshal(v)
}

// Decode deserializes MessagePack bytes into a value.
func (c *MsgPack) Decode(data []byte, target any) error {
	if len(data) == 0 {
		return ErrInvalidData
	}
	return msgpack.Unmarshal(data, target)
}

// Name returns "msgpack".
func (c *MsgPack) Name() string {
	return "msgpack"
}

// Global default codec instance.
var Default Codec = NewMsgPack()
