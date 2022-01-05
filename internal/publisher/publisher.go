package publisher

import (
	"context"
	"encoding/json"
	"time"

	"gbu-scanner/internal/entity"

	"github.com/pkg/errors"
	"github.com/streadway/amqp"
)

// Publisher is implementation for scanner.Publisher interface
type Publisher struct {
	rabbit *amqp.Channel
}

// New returns scanner.Publisher implementation via rabbitmq
func New(rabbit *amqp.Channel) *Publisher {
	return &Publisher{
		rabbit: rabbit,
	}
}

// Init initializes rabbit's entiies like exchanges, queues and queue bindings
func (p *Publisher) Init(ctx context.Context) error {
	err := p.rabbit.ExchangeDeclare(postsExchange, amqp.ExchangeFanout, true, false, false, false, nil)
	if err != nil {
		return errors.Wrap(err, "can't declare exchange")
	}
	return nil
}

func (p *Publisher) Publish(ctx context.Context, post entity.Post) error {
	encoded, err := json.Marshal(post)
	if err != nil {
		return errors.Wrap(err, "can't encode post to JSON")
	}

	err = p.rabbit.Publish(postsExchange, "", false, false, amqp.Publishing{
		DeliveryMode: amqp.Persistent,
		Timestamp:    time.Now(),
		Body:         encoded,
	})
	if err != nil {
		return errors.Wrap(err, "can't publish message")
	}

	return nil
}
