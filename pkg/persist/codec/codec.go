// Package codec provides serialization interfaces and implementations.
package codec

import (
	"errors"
)

// Common errors.
var (
	ErrNilValue        = errors.New("cannot encode nil value")
	ErrInvalidData     = errors.New("invalid data format")
	ErrTypeMismatch    = errors.New("type mismatch during decode")
	ErrUnsupportedType = errors.New("unsupported type for encoding")
)

// Codec handles serialization and deserialization of values.
type Codec interface {
	// Encode serializes a value to bytes.
	Encode(v any) ([]byte, error)

	// Decode deserializes bytes into a value.
	// The target must be a pointer.
	Decode(data []byte, target any) error

	// Name returns the codec name (e.g., "msgpack", "json").
	Name() string
}

// BinaryCodec is the interface for types that can marshal themselves.
type BinaryCodec interface {
	MarshalBinary() ([]byte, error)
	UnmarshalBinary(data []byte) error
}
