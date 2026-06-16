package data

import (
	"context"
	"encoding/json"

	"cascade-oj/pkg/mq"

	amqp "github.com/rabbitmq/amqp091-go"
)

// produce a message to mq to invalidate user contest cache
func (d *Data) publishContestInvalidation(ctx context.Context, msg *mq.ContestCacheMsg) error {
	conn, err := d.mq_channel.GetConnection()
	if err != nil {
		return err
	}
	ch, err := conn.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	q, err := ch.QueueDeclare(
		mq.ContestCacheQueueName, // name
		true,                     // durable
		false,                    // delete when unused
		false,                    // exclusive
		false,                    // no-wait
		nil,                      // arguments
	)
	if err != nil {
		return err
	}

	jsonBody, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	err = ch.PublishWithContext(
		ctx,
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        jsonBody,
		},
	)
	if err != nil {
		return err
	}

	return nil
}

// produce a message to mq to invalidate user problem cache
func (d *Data) publishProblemInvalidation(ctx context.Context, msg *mq.ProblemCacheMsg) error {
	conn, err := d.mq_channel.GetConnection()
	if err != nil {
		return err
	}
	ch, err := conn.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	q, err := ch.QueueDeclare(
		mq.ProblemCacheQueueName, // name
		true,                     // durable
		false,                    // delete when unused
		false,                    // exclusive
		false,                    // no-wait
		nil,                      // arguments
	)
	if err != nil {
		return err
	}

	jsonBody, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	err = ch.PublishWithContext(
		ctx,
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        jsonBody,
		},
	)
	if err != nil {
		return err
	}

	return nil
}
