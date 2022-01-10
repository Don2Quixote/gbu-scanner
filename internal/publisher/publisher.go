package publisher

import (
	"context"
	"encoding/json"
	"sync"
	"time"

	"gbu-scanner/internal/entity"

	"gbu-scanner/pkg/logger"
	"gbu-scanner/pkg/wrappers/rabbit"

	"github.com/pkg/errors"
	"github.com/streadway/amqp"
)

// Publisher is implementation for scanner.Publisher interface.
type Publisher struct {
	rabbitConfig RabbitConfig
	rabbit       *amqp.Channel
	log          logger.Logger
	// RWMutex Locks used to connect to rabbit (Init method)
	// RWMutex RLocks used to use connection
	mu *sync.RWMutex
}

// New returns scanner.Publisher implementation via rabbitmq.
func New(rabbitConfig RabbitConfig, log logger.Logger) *Publisher {
	return &Publisher{
		rabbitConfig: rabbitConfig,
		rabbit:       nil, // Initialized in Init method
		log:          log,
		mu:           &sync.RWMutex{},
	}
}

// Init connects to rabbit and gets rabbit channel, after what
// initializes rabbit's entiies like exchanges, queues etc.
// It also registers a handler for channel closed event to reconnect.
func (p *Publisher) Init(ctx context.Context) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	cfg := p.rabbitConfig
	conn, err := rabbit.Dial(cfg.Host, cfg.User, cfg.Pass, cfg.Vhost, cfg.Amqps)
	if err != nil {
		return errors.Wrap(err, "can't connect to rabbit")
	}

	ch, err := conn.Channel()
	if err != nil {
		return errors.Wrap(err, "can't get rabbit channel")
	}

	err = ch.ExchangeDeclare(postsExchange, amqp.ExchangeFanout, true, false, false, false, nil)
	if err != nil {
		return errors.Wrap(err, "can't declare exchange")
	}

	errs := make(chan *amqp.Error)
	ch.NotifyClose(errs)

	handleChannelClose := func() {
		closeErr := <-errs // This chan will get a value when rabbit channel will be closed
		p.log.Error(errors.Wrap(closeErr, "rabbit channel closed"))

		if !conn.IsClosed() {
			err := conn.Close()
			if err != nil {
				p.log.Error(errors.Wrap(err, "can't close rabbit connection"))
			}
		}

		for attempt, isConnected := 1, false; !isConnected; attempt++ {
			time.Sleep(cfg.ReconnectDelay)

			err := p.Init(ctx)
			if err != nil {
				p.log.Warn(errors.Wrapf(err, "can't re-init publisher (attempt #%d)", attempt))
				continue
			}

			isConnected = true
		}

		p.log.Info("reconnected to rabbit")
	}
	go handleChannelClose()

	p.rabbit = ch

	return nil
}

func (p *Publisher) Publish(ctx context.Context, post entity.Post) error {
	p.mu.RLock()
	defer p.mu.RUnlock()

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
