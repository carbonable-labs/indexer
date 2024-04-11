package dispatcher

import (
	"context"
	"os"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
)

type NatsDispatcherOptsFunc func(*NatsDispatcherOpts)

type NatsDispatcherOpts struct {
	streamConfig jetstream.StreamConfig
	url          string
	bucket       string
	token        string
	timeout      time.Duration
}

func defaultNatsOptions() *NatsDispatcherOpts {
	return &NatsDispatcherOpts{
		url:     nats.DefaultURL,
		bucket:  "storage",
		timeout: 10 * time.Second,
		streamConfig: jetstream.StreamConfig{
			Name:        "EVENTS",
			Subjects:    []string{"*.event.>"},
			Description: "Event stream",
		},
		token: os.Getenv("NATS_TOKEN"),
	}
}

func WithBucket(b string) NatsDispatcherOptsFunc {
	return func(o *NatsDispatcherOpts) {
		o.bucket = b
	}
}

func WithToken(b string) NatsDispatcherOptsFunc {
	return func(o *NatsDispatcherOpts) {
		o.token = b
	}
}

func WithUrl(u string) NatsDispatcherOptsFunc {
	return func(o *NatsDispatcherOpts) {
		if u == "" {
			o.url = nats.DefaultURL
		}
		o.url = u
	}
}

func WithStreamConfig(c jetstream.StreamConfig) NatsDispatcherOptsFunc {
	return func(o *NatsDispatcherOpts) {
		o.streamConfig = c
	}
}

type NatsDispatcher struct {
	js jetstream.JetStream
}

func (nd *NatsDispatcher) Publish(topic string, data []byte) error {
	_, err := nd.js.PublishAsync(topic, data)
	return err
}

func (nd *NatsDispatcher) Subscribe(cName string, topic string, handler MsgHandler) error {
	c, err := nd.js.CreateOrUpdateConsumer(context.Background(), topic, jetstream.ConsumerConfig{
		Name:    cName,
		Durable: cName,
	})
	if err != nil {
		return err
	}

	_, err = c.Consume(func(m jetstream.Msg) {
		handler(m.Data())
	})

	return err
}

func NewNatsDispatcher(opts ...NatsDispatcherOptsFunc) *NatsDispatcher {
	o := defaultNatsOptions()
	for _, optFn := range opts {
		optFn(o)
	}

	nc, _ := nats.Connect(o.url, nats.Token(o.token))

	js, _ := jetstream.New(nc)
	js.CreateOrUpdateStream(context.Background(), o.streamConfig)

	return &NatsDispatcher{js: js}
}
