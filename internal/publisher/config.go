package publisher

import "time"

// RabbitConfig is configuration for rabbitmq's connection
type RabbitConfig struct {
	// Host of rabbitmq
	Host string
	// User for rabbitmq
	User string
	// Pass is password for rabbitmq
	Pass string
	// Vhost is vhost in rabbitmq to connect to
	Vhost string
	// Amqps flag shows should amqps protocol be used instead of amqp or not
	Amqps bool
	// ReconnectDelay is duration how long should wait before
	// attempting to reconnect to rabbit after loosing connection
	ReconnectDelay time.Duration
}
