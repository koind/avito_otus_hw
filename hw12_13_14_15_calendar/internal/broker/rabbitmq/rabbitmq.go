package rabbitmq

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/koind/avito_otus_hw/hw12_13_14_15_calendar/internal/app"
	"github.com/koind/avito_otus_hw/hw12_13_14_15_calendar/internal/domain/entity"
	amqp "github.com/rabbitmq/amqp091-go"
)

type Interface interface {
	GetChannel() *amqp.Channel
	Add(n entity.Notification) error
	GetNotificationChannel() (<-chan entity.Notification, error)
}

type rabbit struct {
	exchange    string
	queue       string
	consumerTag string
	channel     *amqp.Channel
	logger      app.Logger
}

func (q *rabbit) GetChannel() *amqp.Channel {
	return q.channel
}

func New(
	ctx context.Context,
	dsn string,
	exchange string,
	queue string,
	logger app.Logger,
) (Interface, error) {
	conn, err := amqp.Dial(dsn)
	if err != nil {
		return nil, fmt.Errorf("error connect to rabbit (%s): %w", dsn, err)
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("error open rabbit channel (%s): %w", dsn, err)
	}

	if len(exchange) > 0 {
		if err = ch.ExchangeDeclare(
			exchange,
			amqp.ExchangeDirect,
			true,
			false,
			false,
			false,
			nil,
		); err != nil {
			return nil, fmt.Errorf("error declare exchange (%s): %w", exchange, err)
		}
	}

	q, err := ch.QueueDeclare(
		queue,
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("error declare queue (%s): %w", queue, err)
	}

	if err = ch.QueueBind(
		q.Name,
		q.Name,
		exchange,
		false,
		nil,
	); err != nil {
		return nil, fmt.Errorf("error bind queue: %w", err)
	}

	go func() {
		<-ctx.Done()
		ch.Close()
		conn.Close()
	}()

	return &rabbit{
		exchange:    exchange,
		queue:       queue,
		consumerTag: "calendar-consumer",
		channel:     ch,
		logger:      logger,
	}, nil
}

func (q *rabbit) Add(n entity.Notification) error {
	body, err := json.Marshal(n)
	if err != nil {
		return fmt.Errorf("error marshall notification: %w", err)
	}

	if err = q.channel.Publish(
		q.exchange,
		q.queue,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		}); err != nil {
		return fmt.Errorf("error publish notification: %w", err)
	}

	return nil
}

func (q *rabbit) GetNotificationChannel() (<-chan entity.Notification, error) {
	ch := make(chan entity.Notification)

	deliveries, err := q.channel.Consume(
		q.queue,
		q.consumerTag,
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("error consume queue (%s): %w", q.queue, err)
	}

	go func() {
		for d := range deliveries {
			var notification entity.Notification
			if err := json.Unmarshal(d.Body, &notification); err != nil {
				q.logger.Error("error unmarshal notification: %s", err)
				continue
			}

			ch <- notification

			d.Ack(false)
		}

		close(ch)
	}()

	return ch, nil
}
