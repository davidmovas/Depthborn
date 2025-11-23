package entity

type BaseCloneable struct{}

func NewBaseCloneable() *BaseCloneable {
	return &BaseCloneable{}
}

func (bc *BaseCloneable) Clone() any {
	// TODO: Implement proper deep cloning
	return nil
}
