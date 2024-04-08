package dispatcher

import (
	"context"
	"log"
	"os"
	"os/signal"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

const (
	DefaultPublishTimeout = 5 * time.Second
)

type RabbitMQDispatcher struct {
	server   string
	exchange string
	conn     *amqp.Connection
	ch       *amqp.Channel
}

func NewRabbitMQDispatcher(server string, exchange string) *RabbitMQDispatcher {
	rd := &RabbitMQDispatcher{
		server:   server,
		exchange: exchange,
	}

	conn, err := amqp.Dial(rd.server)
	failOnError(err, "Failed to connect to RabbitMQ")

	rd.conn = conn

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")

	rd.ch = ch

	err = ch.ExchangeDeclare(
		rd.exchange, // exchange name
		"direct",    // type
		true,        // durable
		false,       // auto-deleted
		false,       // internal
		false,       // no-wait
		nil,         // arguments
	)
	failOnError(err, "Failed to declare an exchange")

	return rd
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%v: %v", msg, err)
	}
}

func (rd *RabbitMQDispatcher) Publish(topic string, data []byte) error {
	ctx, cancelFn := context.WithTimeout(context.Background(), DefaultPublishTimeout)
	defer cancelFn()

	if err := rd.ch.PublishWithContext(ctx,
		rd.exchange, // exchange
		topic,       // routing key
		false,       // mandatory
		false,       // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        data,
		}); err != nil {
		return err
	}

	return nil
}

func (rd *RabbitMQDispatcher) Subscribe(cName string, topic string, handler MsgHandler) error {
	// Graceful shutdown the subscriber
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill)
	defer stop()

	q, err := rd.ch.QueueDeclare(
		cName,
		false,
		false,
		true,
		false,
		nil,
	)
	failOnError(err, "Failed to declare a queue")

	err = rd.ch.QueueBind(
		q.Name,
		topic,
		rd.exchange,
		false,
		nil,
	)
	failOnError(err, "Failed to bind a queue")

	msgs, err := rd.ch.Consume(
		q.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	go func() {
		for d := range msgs {
			if err := handler(d.Body); err != nil {
				log.Printf("Failed to handle message: %v", err)
			}
		}
	}()

	log.Printf(" [*] Processing messages. To exit press CTRL+C")

	for {
		select {
		case <-ctx.Done():
			log.Println("Shutting down the subscriber")
			rd.Close()
			return nil
		}
	}
}

func (rd *RabbitMQDispatcher) Close() {
	rd.ch.Close()
	rd.conn.Close()
}
