package main

import (
	"os"

	"github.com/hashicorp/go-hclog"
	"github.com/streadway/amqp"
)

var rabbit_user = os.Getenv("RABBIT_USERNAME")
var rabbit_password = os.Getenv("RABBIT_PASSWORD")
var rabbit_port = os.Getenv("RABBIT_PORT")
var rabbit_host = os.Getenv("RABBIT_HOST")

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
	url:= "amqp://" + rabbit_user + ":" + rabbit_password + "@" + rabbit_host + ":" + rabbit_port + "/"
	println(url)
	conn, err := amqp.Dial(url)
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
