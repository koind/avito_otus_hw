package mq

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/koind/avito_otus_hw/hw12_13_14_15_calendar/internal/app"
	amqp "github.com/rabbitmq/amqp091-go"
)

type Rabbit struct {
	exchange    string
	queue       string
	consumerTag string
	channel     *amqp.Channel
	logger      app.Logger
}

func NewRabbit(
	ctx context.Context,
	dsn string,
	exchange string,
	queue string,
	logger app.Logger,
) (*Rabbit, error) {
	conn, err := amqp.Dial(dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ on %s: %w", dsn, err)
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("failed to open RabbitMQ Channel on %s: %w", dsn, err)
	}

	if len(exchange) > 0 {
		err = ch.ExchangeDeclare(
			exchange,
			amqp.ExchangeDirect,
			true,
			false,
			false,
			false,
			nil,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to declare an exchanhe %s: %w", exchange, err)
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
		return nil, fmt.Errorf("failed to declare a queue %s: %w", queue, err)
	}

	err = ch.QueueBind(
		q.Name,
		q.Name,
		exchange,
		false,
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to bind queue: %w", err)
	}

	go func() {
		<-ctx.Done()
		ch.Close()
		conn.Close()
	}()

	return &Rabbit{
		exchange:    exchange,
		queue:       queue,
		consumerTag: "calendar-consumer",
		channel:     &amqp.Channel{},
		logger:      logger,
	}, nil
}

func (q *Rabbit) Add(n app.Notification) error {
	body, err := json.Marshal(n)
	if err != nil {
		return fmt.Errorf("failed to marshall notification: %w", err)
	}

	err = q.channel.Publish(
		q.exchange,
		q.queue,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		})

	if err != nil {
		return fmt.Errorf("failed to publish notification: %w", err)
	}

	return nil
}

func (q *Rabbit) GetNotificationChannel() (<-chan app.Notification, error) {
	ch := make(chan app.Notification)

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
		return nil, fmt.Errorf("failed to consume queue %s: %w", q.queue, err)
	}

	go func() {
		for d := range deliveries {
			var notification app.Notification
			err := json.Unmarshal(d.Body, &notification)
			if err != nil {
				q.logger.Error("Failed to unmarshal notification message: %s", err)
				continue
			}

			ch <- notification

			d.Ack(false)
		}

		close(ch)
	}()

	return ch, nil
}
