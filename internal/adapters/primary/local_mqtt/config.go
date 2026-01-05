package local_mqtt

import (
	"os"
	"strconv"
	"time"
)

type Config struct {
	Broker      string
	Port        int
	ClientID    string
	Username    string
	Password    string
	QoS         byte
	KeepAlive   time.Duration
	PingTimeout time.Duration
}

func NewConfig() Config {
	return Config{
		Broker:      getenvDefault("MQTT_BROKER", "tcp://localhost"),
		Port:        getenvIntDefault("MQTT_PORT", 1883),
		ClientID:    getenvDefault("MQTT_CLIENT_ID", "server-mqtt"),
		Username:    os.Getenv("MQTT_USERNAME"),
		Password:    os.Getenv("MQTT_PASSWORD"),
		QoS:         byte(getenvIntDefault("MQTT_QOS", 1)),
		KeepAlive:   getenvDurationDefault("MQTT_KEEPALIVE", 30*time.Second),
		PingTimeout: getenvDurationDefault("MQTT_PING_TIMEOUT", 10*time.Second),
	}
}

func (c Config) BrokerURL() string {
	return c.Broker + ":" + strconv.Itoa(c.Port)
}

func getenvDefault(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

func getenvIntDefault(key string, fallback int) int {
	if value := os.Getenv(key); value != "" {
		if parsed, err := strconv.Atoi(value); err == nil {
			return parsed
		}
	}
	return fallback
}

func getenvDurationDefault(key string, fallback time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		if parsed, err := time.ParseDuration(value); err == nil {
			return parsed
		}
	}
	return fallback
}
