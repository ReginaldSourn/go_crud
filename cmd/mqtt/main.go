package main

import (
	"context"
	"log"
	"os/signal"
	"syscall"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/joho/godotenv"
	localmqtt "github.com/reginaldsourn/go-crud/internal/adapters/primary/local_mqtt"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("no .env file found; using existing environment variables")
	}

	cfg := localmqtt.NewConfig()

	client, err := localmqtt.NewClient(cfg)
	if err != nil {
		log.Fatalf("failed to create mqtt client: %v", err)
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	connectCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	if err := client.Connect(connectCtx); err != nil {
		log.Fatalf("mqtt connect failed: %v", err)
	}

	if err := registerSubscriptions(client); err != nil {
		log.Fatalf("failed to register subscriptions: %v", err)
	}

	log.Println("mqtt service started, waiting for messages...")
	<-ctx.Done()

	log.Println("mqtt service shutting down")
	client.Disconnect(250)
	log.Println("mqtt service stopped")
}

func registerSubscriptions(client *localmqtt.Client) error {
	topics := []string{
		"devices/+/status",
		"devices/+/telemetry",
	}

	for _, topic := range topics {
		if err := client.Subscribe(topic, createMessageHandler(topic)); err != nil {
			return err
		}
	}

	return nil
}

func createMessageHandler(topic string) mqtt.MessageHandler {
	return func(client mqtt.Client, msg mqtt.Message) {
		log.Printf("received message: topic=%s payload=%s", msg.Topic(), string(msg.Payload()))
		handleMessage(msg.Topic(), msg.Payload())
	}
}

func handleMessage(topic string, payload []byte) {
	// TODO: Implement message handling logic based on topic patterns
	// Example:
	// - devices/+/status -> update device status in database
	// - devices/+/telemetry -> store telemetry data
}
