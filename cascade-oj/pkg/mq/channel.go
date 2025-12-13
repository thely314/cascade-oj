package mq

import (
	"context"
	"encoding/json"

	"github.com/go-kratos/kratos/v2/log"
	amqp "github.com/rabbitmq/amqp091-go"
)

// specific queue names for go-judge engine
const (
	GojudgeSubmissionQueueName = "gojudge-submission-queue"
	GojudgeSelfTestQueueName   = "gojudge-self-test-queue"
)

type Channel interface {
	Publish(ctx context.Context, msg interface{}) error
	Close() error
}

type CascadeOjQueue struct {
	ojCh  *amqp.Channel
	queue *amqp.Queue
}

type CascadeOjExchange struct {
	ojCh  *amqp.Channel
	topic string
}

// creates a new queue with the given name
func NewCascadeOjQueue(ch *amqp.Channel, queueName string) (*CascadeOjQueue, error) {
	q, err := ch.QueueDeclare(
		queueName, // name
		true,      // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	if err != nil {
		return nil, err
	}
	return &CascadeOjQueue{
		ojCh:  ch,
		queue: &q,
	}, nil
}

// creates a new exchange with the given name
func NewCascadeOjExchange(ch *amqp.Channel, exchangeName string) (*CascadeOjExchange, error) {
	err := ch.ExchangeDeclare(
		exchangeName+"_topic", // name
		"topic",               // type
		true,                  // durable
		false,                 // auto-deleted
		false,                 // internal
		false,                 // no-wait
		nil,                   // arguments
	)
	if err != nil {
		return nil, err
	}
	return &CascadeOjExchange{
		ojCh:  ch,
		topic: exchangeName,
	}, nil
}

// publishes a message to the queue
func (casq *CascadeOjQueue) Publish(ctx context.Context, msg interface{}) error {
	body, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	log.Infof("Publishing message to queue %s: %s", casq.queue.Name, string(body))
	return casq.ojCh.PublishWithContext(ctx,
		"",              // exchange
		casq.queue.Name, // routing key
		false,           // mandatory
		false,           // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
}

func (casq *CascadeOjQueue) Close() error {
	return casq.ojCh.Close()
}

// publishes a message to the exchange
func (casex *CascadeOjExchange) Publish(ctx context.Context, msg interface{}) error {
	body, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	log.Infof("Publishing message to exchange %s: %s", casex.topic, string(body))
	return casex.ojCh.PublishWithContext(ctx,
		casex.topic, // exchange
		casex.topic, // routing key
		false,       // mandatory
		false,       // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
}

func (casex *CascadeOjExchange) Close() error {
	return casex.ojCh.Close()
}
