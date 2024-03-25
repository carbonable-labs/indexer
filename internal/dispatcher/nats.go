package dispatcher

import (
	"log"

	"github.com/nats-io/nats.go"
)

type NatsDispatcher struct {
	connection *nats.Conn
	connError  error
}

func (nd *NatsDispatcher) Publish(topic string, data []byte) error {
	err := nd.connection.Publish(topic, data)
	if err != nil {
		log.Fatal("Error while publishing to topic ", topic, "\nerror: ", err)
		return err
	}
	log.Print("Published to topic [", topic, "] message: '", string(data), "'")
	return nil
}

func (nd *NatsDispatcher) Subscribe(topic string, handler nats.MsgHandler) error {
	nd.connection.Subscribe(topic, handler)
	log.Printf("Listening on [%s]", topic)
	return nil
}

func (nd *NatsDispatcher) Connect(server string) {
	nd.connection, nd.connError = nats.Connect(server)
	if nd.connError != nil {
		log.Fatal(nd.connError)
	}
}
