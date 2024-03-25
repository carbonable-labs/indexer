package dispatcher

import (
	"log"
	"testing"

	"github.com/nats-io/nats.go"
)

func TestNatsDispatcher(t *testing.T) {
	// to run this test you must use the following docker command
	// docker run -p 4222:4222 -ti nats:latest
	//and make sure the docker container for nats is running

	server := "nats://localhost:4222"
	subject := "Costa Rica weather"
	message := "Sunny with possibility of rain"

	natsPublisher := NatsDispatcher{}
	natsPublisher.Connect(server)

	natsSubscriber := NatsDispatcher{}
	natsSubscriber.Connect(server)

	natsSubscriber.Subscribe(subject, func(msg *nats.Msg) {
		log.Printf("Received on [%s]: '%s'", msg.Subject, string(msg.Data))
	})

	natsPublisher.Publish(subject, []byte(message))
}
