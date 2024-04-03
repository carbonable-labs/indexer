package storage

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/charmbracelet/log"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
)

type NatsStorageOptsFunc func(*NatsStorageOptions)

type NatsStorageOptions struct {
	url     string
	bucket  string
	timeout time.Duration
}

func defaultNatsOptions() *NatsStorageOptions {
	return &NatsStorageOptions{
		url:     nats.DefaultURL,
		bucket:  "storage",
		timeout: 10 * time.Second,
	}
}

func WithBucket(b string) NatsStorageOptsFunc {
	return func(o *NatsStorageOptions) {
		o.bucket = b
	}
}

type NatsStorage struct {
	kv jetstream.KeyValue
}

func (s *NatsStorage) Get(id []byte) []byte {
	entry, err := s.kv.Get(context.Background(), string(id))
	if err != nil {
		log.Error("failed to get value from nats kv", "error", err)
		return []byte{}
	}

	return entry.Value()
}

func (s *NatsStorage) Has(id []byte) bool {
	_, err := s.kv.Get(context.Background(), string(id))

	return errors.Is(err, jetstream.ErrKeyNotFound)
}

func (s *NatsStorage) Set(key []byte, value []byte) error {
	_, err := s.kv.Put(context.Background(), string(key), value)
	return err
}

func (s *NatsStorage) Scan(prefix []byte) [][]byte {
	var res [][]byte
	entries, err := s.kv.Watch(context.Background(), fmt.Sprintf("%s>", string(prefix)), jetstream.IncludeHistory())
	defer entries.Stop()
	if err != nil {
		log.Error("failed to scan value from nats kv", "error", err)
		return [][]byte{}
	}
	for v := range entries.Updates() {
		if v == nil {
			return res
		}
		res = append(res, v.Value())
	}

	return res
}

func NewNatsStorage(opts ...NatsStorageOptsFunc) *NatsStorage {
	o := defaultNatsOptions()
	for _, optFn := range opts {
		optFn(o)
	}

	nc, _ := nats.Connect(o.url)

	js, _ := jetstream.New(nc)

	ctx, cancel := context.WithTimeout(context.Background(), o.timeout)
	defer cancel()

	kv, _ := js.CreateOrUpdateKeyValue(ctx, jetstream.KeyValueConfig{
		Bucket: o.bucket,
	})

	return &NatsStorage{kv: kv}
}
