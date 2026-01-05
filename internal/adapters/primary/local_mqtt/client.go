package local_mqtt

import (
	"context"
	"fmt"
	"log"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type Client struct {
	client mqtt.Client
	config Config
}

func NewClient(cfg Config) (*Client, error) {
	opts := mqtt.NewClientOptions().
		AddBroker(cfg.BrokerURL()).
		SetClientID(cfg.ClientID).
		SetUsername(cfg.Username).
		SetPassword(cfg.Password).
		SetKeepAlive(cfg.KeepAlive).
		SetPingTimeout(cfg.PingTimeout).
		SetAutoReconnect(true).
		SetConnectRetry(true).
		SetConnectRetryInterval(5 * time.Second)

	opts.OnConnect = func(c mqtt.Client) {
		log.Printf("mqtt connected: broker=%s client_id=%s", cfg.BrokerURL(), cfg.ClientID)
	}

	opts.OnConnectionLost = func(c mqtt.Client, err error) {
		log.Printf("mqtt connection lost: %v", err)
	}

	opts.OnReconnecting = func(c mqtt.Client, opts *mqtt.ClientOptions) {
		log.Printf("mqtt reconnecting to %s", cfg.BrokerURL())
	}

	client := mqtt.NewClient(opts)

	return &Client{
		client: client,
		config: cfg,
	}, nil
}

func (c *Client) Connect(ctx context.Context) error {
	log.Printf("mqtt connecting to broker: %s", c.config.BrokerURL())

	token := c.client.Connect()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-token.Done():
		if err := token.Error(); err != nil {
			return fmt.Errorf("mqtt connect failed: %w", err)
		}
	}

	return nil
}

func (c *Client) Disconnect(waitMs uint) {
	log.Println("mqtt disconnecting")
	c.client.Disconnect(waitMs)
}

func (c *Client) Subscribe(topic string, handler mqtt.MessageHandler) error {
	token := c.client.Subscribe(topic, c.config.QoS, handler)
	if ok := token.WaitTimeout(10 * time.Second); !ok {
		return fmt.Errorf("mqtt subscribe timeout for topic: %s", topic)
	}
	if err := token.Error(); err != nil {
		return fmt.Errorf("mqtt subscribe failed: %w", err)
	}
	log.Printf("mqtt subscribed to topic: %s", topic)
	return nil
}

func (c *Client) Publish(topic string, payload interface{}) error {
	token := c.client.Publish(topic, c.config.QoS, false, payload)
	if ok := token.WaitTimeout(10 * time.Second); !ok {
		return fmt.Errorf("mqtt publish timeout for topic: %s", topic)
	}
	if err := token.Error(); err != nil {
		return fmt.Errorf("mqtt publish failed: %w", err)
	}
	return nil
}

func (c *Client) Unsubscribe(topics ...string) error {
	token := c.client.Unsubscribe(topics...)
	if ok := token.WaitTimeout(10 * time.Second); !ok {
		return fmt.Errorf("mqtt unsubscribe timeout")
	}
	if err := token.Error(); err != nil {
		return fmt.Errorf("mqtt unsubscribe failed: %w", err)
	}
	return nil
}

func (c *Client) IsConnected() bool {
	return c.client.IsConnected()
}
