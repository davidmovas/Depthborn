package codec

import (
	"encoding/json"
)

// JSON implements Codec using JSON serialization.
// Useful for debugging and human-readable output.
type JSON struct {
	indent bool
}

// NewJSON creates a new JSON codec.
func NewJSON() *JSON {
	return &JSON{indent: false}
}

// NewJSONIndented creates a new JSON codec with pretty-printing.
func NewJSONIndented() *JSON {
	return &JSON{indent: true}
}

// Encode serializes a value to JSON bytes.
func (c *JSON) Encode(v any) ([]byte, error) {
	if v == nil {
		return nil, ErrNilValue
	}
	if c.indent {
		return json.MarshalIndent(v, "", "  ")
	}
	return json.Marshal(v)
}

// Decode deserializes JSON bytes into a value.
func (c *JSON) Decode(data []byte, target any) error {
	if len(data) == 0 {
		return ErrInvalidData
	}
	return json.Unmarshal(data, target)
}

// Name returns "json".
func (c *JSON) Name() string {
	return "json"
}
