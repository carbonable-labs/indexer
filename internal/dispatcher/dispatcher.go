package dispatcher

type (
	MsgHandler      func(data []byte) error
	EventDispatcher interface {
		Publish(topic string, data []byte) error
		Subscribe(topic string, handler MsgHandler) error
	}
)
