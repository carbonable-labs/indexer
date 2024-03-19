package dispatcher

type InMemoryEventDispatcher struct{}

func (d *InMemoryEventDispatcher) Publish(topic string, data []byte) error {
	panic("TODO")
}

func (d *InMemoryEventDispatcher) Subscribe(topic string, handler MsgHandler) error {
	panic("TODO")
}

func NewInMemoryBus() EventDispatcher {
	return &InMemoryEventDispatcher{}
}
