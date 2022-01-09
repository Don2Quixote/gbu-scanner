package rabbit

import (
	"fmt"

	"github.com/streadway/amqp"
)

// Dial is wrapper for original amqp.Dial function but with
// credential-parameters instead of connection string.
func Dial(host, user, pass, vhost string, amqps bool) (*amqp.Connection, error) {
	connString := fmt.Sprintf("amqp://%s:%s@%s/%s", user, pass, host, vhost)
	if amqps {
		connString = fmt.Sprintf("amqps://%s:%s@%s/%s", user, pass, host, vhost)
	}
	return amqp.Dial(connString)
}
