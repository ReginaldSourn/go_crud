package ports

import "context"

type MQTTClient interface {
	Connect(ctx context.Context) error
	Disconnect(waitMs uint)
	Subscribe(topic string, handler MessageHandler) error
	Publish(topic string, payload interface{}) error
	Unsubscribe(topics ...string) error
	IsConnected() bool
}

type MessageHandler func(topic string, payload []byte)
