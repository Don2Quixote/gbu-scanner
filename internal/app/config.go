package app

import "gbu-scanner/pkg/logger"

// appConfig is struct for parsing ENV configuration
type appConfig struct {
	// BlogHost is host where blog is located (I guess it will always "go.dev")
	// If BlogHost empty - setDefaults method will set BlogHost, BlogPath and BlogHttps
	BlogHost string `config:"BLOG_HOST"`
	// BlogPath is path where blog is located (I guess it will always be "/blog/all")
	BlogPath string `config:"BLOG_PATH"`
	// BlogHTTPS flag shows should https protocol used instead of http or not
	BlogHTTPS bool `config:"BLOG_HTTPS"`
	// BlogScanInterval is delay (in seconds) between blog's scans
	BlogScanInterval int `config:"BLOG_SCAN_INTERVAL,required"`
	// BlogScanNetworkTimeout is http client's timeout (in seconds) during request to blog
	BlogScanNetworkTimeout int `config:"BLOG_SCAN_NETWORK_TIMEOUT,required"`
	// MongoHost is host of mongodb
	MongoHost string `config:"MONGO_HOST,required"`
	// MongoUser is user for mongodb
	MongoUser string `config:"MONGO_USER"`
	// MongoPass is password for mongodb
	MongoPass string `config:"MONGO_PASS"`
	// MongoDatabase is database's name in mongodb
	MongoDatabase string `config:"MONGO_DATABASE"`
	// MongoSRV flag shows should mongodb+srv protocol used instead of just mongo or not
	MongoSRV bool `config:"MONGO_SRV"`
	// RabbitHost is host of rabbitmq
	RabbitHost string `config:"RABBIT_HOST,required"`
	// RabbitUser is user for rabbitmq
	RabbitUser string `config:"RABBIT_USER"`
	// RabbitPass is password for rabbitmq
	RabbitPass string `config:"RABBIT_PASS"`
	// RabbitVhost is vhost in rabbitmq to connect
	RabbitVhost string `config:"RABBIT_VHOST"`
	// RabbitAmqps flag shows should amqps protocol be used instead of amqp or not
	RabbitAmqps bool `config:"RABBIT_AMQPS"`
	// RabbitReconnectDelay is delay (in seconds) before attempting to reconnect to rabbit after loosing connection
	RabbitReconnectDelay int `config:"RABBIT_RECONNECT_DELAY,required"`
}

// setDefaults sets some default config variables if they are empty
func (c *appConfig) setDefaults(log logger.Logger) {
	if c.BlogHost == "" {
		log.Warn("BlogHost config var is empty, setting BlogHost, BlogPath and BlogHttps to defaults")
		c.BlogHost = "go.dev"
		c.BlogPath = "/blog/all"
		c.BlogHTTPS = true
	}
}
