# dispatcher

Public interface in `dispatcher.go`

- `Publish(topic string, data []byte)`
- `Subscribe(topic string, handler func(data []byte))`

Implementations for `nats` & `rabbitmq` present in files with the same name.
