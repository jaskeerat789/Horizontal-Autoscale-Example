package main

import (
	"github.com/hashicorp/go-hclog"
	"github.com/streadway/amqp"
)

type RabbitMQClient struct {
	l    hclog.Logger
	conn *amqp.Connection
	ch   *amqp.Channel
	q    amqp.Queue
}

func NewClient() *RabbitMQClient {
	log := hclog.New(&hclog.LoggerOptions{
		Name: "RabbitMQ",
	})

	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(log, err, "Error establishing connection with RabbitMQ")

	ch, err := conn.Channel()
	failOnError(log, err, "Failed to create a channel")

	q, _ := creatQueue("Orders", log, ch)

	return &RabbitMQClient{
		l:    log,
		conn: conn,
		ch:   ch,
		q:    q,
	}
}

func creatQueue(name string, l hclog.Logger, ch *amqp.Channel) (amqp.Queue, error) {
	l.Info("Creating queue", "name", name)
	q, err := ch.QueueDeclare(
		name,
		false,
		false,
		false,
		false,
		nil,
	)
	failOnError(l, err, "Failed to declare Queue: Name"+name)
	return q, err
}

func (rc *RabbitMQClient) SendMessage(body []byte) error {
	rc.l.Info("Publishing message", "queue", rc.q.Name, "message", string(body[:]))
	err := rc.ch.Publish(
		"",
		rc.q.Name,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        body,
		},
	)
	failOnError(rc.l, err, "Failed to publish a message")
	return err
}

func failOnError(log hclog.Logger, err error, msg string) {
	if err != nil {
		log.Error(msg, "error", err)
	}
}
