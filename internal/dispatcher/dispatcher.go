package dispatcher

type (
	MsgHandler      func(data []byte) error
	EventDispatcher interface {
		Publish(topic string, data []byte) error
		Subscribe(cName string, topic string, handler MsgHandler) error
	}
)
