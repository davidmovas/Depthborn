package serializer

import "github.com/vmihailenco/msgpack/v5"

var _ Serializer = (*MessagePackSerializer)(nil)

type MessagePackSerializer struct{}

func NewMessagePackSerializer() Serializer {
	return &MessagePackSerializer{}
}

func (s *MessagePackSerializer) Marshal(v any) ([]byte, error) {
	return msgpack.Marshal(v)
}

func (s *MessagePackSerializer) Unmarshal(data []byte, v any) error {
	return msgpack.Unmarshal(data, v)
}

func (s *MessagePackSerializer) Name() string {
	return "msgpack"
}
