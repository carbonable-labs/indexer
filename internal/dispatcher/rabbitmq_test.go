package dispatcher

import (
	"log"
	"syscall"
	"testing"
	"time"

	"github.com/test-go/testify/assert"
)

func TestRabbitMQDispatcher(t *testing.T) {
	// to run this test you must use the following docker command
	// docker run -p 5672:5672 -ti rabbitmq:latest
	// and make sure the docker container for RabbitMQ is running by checking `docker ps`

	server := "amqp://guest:guest@localhost:5672/"
	exchangeName := "test-exchange"
	topic := "weather"
	message := "Sunny with possibility of rain"
	message2 := "Thunder storm"

	rabbitmqPublisher := NewRabbitMQDispatcher(server, exchangeName)
	rabbitmqSubscriber := NewRabbitMQDispatcher(server, exchangeName)

	go rabbitmqSubscriber.Subscribe(topic, func(data []byte) error {
		log.Printf("Received on exchange: [%s], topic: [%s]: '%s'", rabbitmqPublisher.exchange, topic, string(data))

		return nil
	})

	// Sleep 1 sec to make sure the subscriber is ready
	time.Sleep(1 * time.Second)

	err := rabbitmqPublisher.Publish(topic, []byte(message))
	assert.Nil(t, err)

	err = rabbitmqPublisher.Publish(topic, []byte(message2))
	assert.Nil(t, err)

	time.Sleep(2 * time.Second)

	// Send SIGINT to stop the subscriber goroutine
	syscall.Kill(syscall.Getpid(), syscall.SIGINT)

	time.Sleep(1 * time.Second)
}
