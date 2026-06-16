package mq

import (
	amqp "github.com/rabbitmq/amqp091-go"
)

// specific queue names for go-judge engine
const (
	GojudgeSubmissionQueueName = "gojudge-submission-queue"
	GojudgeSelfTestQueueName   = "gojudge-self-test-queue"
)

// specific queue names
const (
	ContestCacheQueueName = "contest-cache-queue"
	ProblemCacheQueueName = "problem-cache-queue"
)

// used for rabbitmq exchange definitions
const (
	ContestExchangeName = "contest_ex_topic"
	ProblemExchangeName = "problem_ex_topic"
	JudgeExchangeName   = "judge_ex_topic"
)

// specific msg types
type ContestCacheMsg struct {
	ContestID int64  `json:"contest_id"`
	Scale     string `json:"scale"` // 缓存级别, "list" or "single"
}

type ProblemCacheMsg struct {
	ProblemID int64  `json:"problem_id"` // single -> problemID, list -> contestID
	Scale     string `json:"scale"`      // 缓存级别, "list" or "single"
}

func NewExchangeDeclare(ch *amqp.Channel, exchangeName, exType string) error {
	return ch.ExchangeDeclare(
		exchangeName, // name
		exType,       // type
		true,         // durable
		false,        // auto-delete
		false,        // internal
		false,        // no-wait
		nil,          // arguments
	)
}

func NewQueueDeclare(ch *amqp.Channel, queueName string) (amqp.Queue, error) {
	return ch.QueueDeclare(
		queueName, // name
		true,      // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
}

func NewBindingDeclare(ch *amqp.Channel, queueName, routingKey, exchangeName string) error {
	return ch.QueueBind(
		queueName,
		routingKey,
		exchangeName,
		false, // no-wait
		nil,   // arguments
	)
}
