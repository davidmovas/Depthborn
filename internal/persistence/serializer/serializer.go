package serializer

// Serializer provides unified interface for data serialization
type Serializer interface {
	// Name returns serializer name for debugging
	Name() string

	// Marshal encodes value into bytes
	Marshal(v any) ([]byte, error)

	// Unmarshal decodes bytes into value
	Unmarshal(data []byte, v any) error
}
