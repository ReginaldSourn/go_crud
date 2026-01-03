package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("no .env file found; using existing environment variables")
	}

	broker := getenvDefault("MQTT_BROKER", "tcp://localhost:1883")
	clientID := getenvDefault("MQTT_CLIENT_ID", "server-mqtt")
	username := os.Getenv("MQTT_USERNAME")
	password := os.Getenv("MQTT_PASSWORD")
	log.Println("mqtt connecting to broker:", broker)
	opts := mqtt.NewClientOptions().
		AddBroker(broker).
		SetClientID(clientID).
		SetUsername(username).
		SetPassword(password).
		SetKeepAlive(30 * time.Second).
		SetPingTimeout(10 * time.Second).
		SetAutoReconnect(true).
		SetConnectRetry(true).
		SetConnectRetryInterval(5 * time.Second)

	opts.OnConnect = func(c mqtt.Client) {
		log.Printf("mqtt connected: broker=%s client_id=%s", broker, clientID)
	}
	opts.OnConnectionLost = func(c mqtt.Client, err error) {
		log.Printf("mqtt connection lost: %v", err)
	}

	client := mqtt.NewClient(opts)
	token := client.Connect()
	if ok := token.WaitTimeout(10 * time.Second); !ok {
		log.Fatal("mqtt connect timeout")
	}
	if err := token.Error(); err != nil {
		log.Fatalf("mqtt connect failed: %v", err)
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	<-ctx.Done()
	log.Println("mqtt shutting down")

	client.Disconnect(250)
}

func getenvDefault(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
