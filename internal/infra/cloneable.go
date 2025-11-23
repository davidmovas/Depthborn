package infra

var _ Cloneable = (*BaseCloneable)(nil)

// BaseCloneable provides common cloning functionality
type BaseCloneable struct{}

func (bc *BaseCloneable) Clone() any {
	// TODO: Implement deep cloning logic
	return nil
}
