package dispatcher

import (
	"github.com/nats-io/nats.go"
)

type NATSEventDispatcher struct {
	conn *nats.Conn
}

// The NewNATSEventDispatcher function creates a new NATSEventDispatcher with a connected NATS client.
func NewNATSEventDispatcher(natsURL string) (*NATSEventDispatcher, error) {
	conn, err := nats.Connect(natsURL)
	if err != nil {
		return nil, err
	}
	return &NATSEventDispatcher{
		conn: conn,
	}, nil
}

// Publish sends a message to a NATS subject.
func (n *NATSEventDispatcher) Publish(topic string, data []byte) error {
	return n.conn.Publish(topic, data)
}

// Subscribe listens for messages on a NATS subject and handles them using the provided MsgHandler.
func (n *NATSEventDispatcher) Subscribe(topic string, handler MsgHandler) error {
	_, err := n.conn.Subscribe(topic, func(msg *nats.Msg) {
		err := handler(msg.Data)
		if err != nil {
			// Handle the error according to your application's requirements.
			// Logging the error is a common approach.
		}
	})
	return err
}
