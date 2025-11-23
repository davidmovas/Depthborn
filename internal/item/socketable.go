package item

type BaseSocketable struct {
	*BaseItem
	socketType SocketType
	effect     SocketEffect
}

func NewBaseSocketable(id string, itemType Type, name string, socketType SocketType) *BaseSocketable {
	return &BaseSocketable{
		BaseItem:   NewBaseItem(id, itemType, name),
		socketType: socketType,
		effect:     nil,
	}
}

func (bs *BaseSocketable) SocketType() SocketType {
	return bs.socketType
}

func (bs *BaseSocketable) Effect() SocketEffect {
	return bs.effect
}
