package storage

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/charmbracelet/log"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
)

type NatsStorageOptsFunc func(*NatsStorageOptions)

type NatsStorageOptions struct {
	url     string
	bucket  string
	token   string
	timeout time.Duration
}

func defaultNatsOptions() *NatsStorageOptions {
	return &NatsStorageOptions{
		url:     nats.DefaultURL,
		bucket:  "storage",
		timeout: 10 * time.Second,
		token:   os.Getenv("NATS_TOKEN"),
	}
}

func WithBucket(b string) NatsStorageOptsFunc {
	return func(o *NatsStorageOptions) {
		o.bucket = b
	}
}

func WithToken(b string) NatsStorageOptsFunc {
	return func(o *NatsStorageOptions) {
		o.token = b
	}
}

func WithUrl(u string) NatsStorageOptsFunc {
	return func(o *NatsStorageOptions) {
		if u == "" {
			o.url = nats.DefaultURL
			return
		}
		o.url = u
	}
}

type NatsStorage struct {
	kv jetstream.KeyValue
}

func (s *NatsStorage) Get(id []byte) []byte {
	ctx, _ := s.createCtx()
	entry, err := s.kv.Get(ctx, string(id))
	if err != nil {
		log.Error("failed to get value from nats kv", "error", err)
		return []byte{}
	}

	return entry.Value()
}

func (s *NatsStorage) Has(id []byte) bool {
	ctx, _ := s.createCtx()
	_, err := s.kv.Get(ctx, string(id))

	return !errors.Is(err, jetstream.ErrKeyNotFound)
}

func (s *NatsStorage) Set(key []byte, value []byte) error {
	ctx, _ := s.createCtx()
	_, err := s.kv.Put(ctx, string(key), value)
	return err
}

// NOTE: nats scan may take quite some time to execute because of the number of keys.
// nats add a default timeout of 5 seconds if we do not provide any timeout.
func (s *NatsStorage) Scan(prefix []byte) [][]byte {
	var res [][]byte
	ctx, _ := s.createCtx()
	entries, err := s.kv.Watch(ctx, fmt.Sprintf("%s>", string(prefix)), jetstream.IncludeHistory())
	if err != nil {
		log.Error("failed to scan value from nats kv", "error", err, "prefix", string(prefix))
		return [][]byte{}
	}
	defer entries.Stop()

	for v := range entries.Updates() {
		if v == nil {
			return res
		}
		res = append(res, v.Value())
	}

	return res
}

func (s *NatsStorage) createCtx() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), 2*time.Minute)
}

func NewNatsStorage(opts ...NatsStorageOptsFunc) *NatsStorage {
	o := defaultNatsOptions()
	for _, optFn := range opts {
		optFn(o)
	}

	natsOpts := []nats.Option{}
	if o.token != "" {
		natsOpts = append(natsOpts, nats.Token(o.token))
	}
	nc, err := nats.Connect(o.url, natsOpts...)
	if err != nil {
		log.Error("failed to connect to nats storage", "error", err)
	}

	js, _ := jetstream.New(nc)

	ns := &NatsStorage{}
	ctx, _ := ns.createCtx()
	kv, err := js.CreateOrUpdateKeyValue(ctx, jetstream.KeyValueConfig{
		Bucket: o.bucket,
	})
	if err != nil {
		log.Error("failed to create or update key value", "error", err)
	}

	ns.kv = kv
	return ns
}
