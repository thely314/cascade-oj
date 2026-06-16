package mq

import (
	"sync"

	amqp "github.com/rabbitmq/amqp091-go"
)

type MQConnection struct {
	url  string
	conn *amqp.Connection
	rwmu sync.RWMutex
}

func NewMQConnection(url string) (*MQConnection, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, err
	}
	return &MQConnection{
		url:  url,
		conn: conn,
	}, nil
}

func (mq *MQConnection) GetConnection() (*amqp.Connection, error) {
	mq.rwmu.RLock()
	if mq.conn != nil && !mq.conn.IsClosed() {
		mq.rwmu.RUnlock()
		return mq.conn, nil
	}
	mq.rwmu.RUnlock()
	return mq.reconnect()
}

func (mq *MQConnection) reconnect() (*amqp.Connection, error) {
	mq.rwmu.Lock()
	defer mq.rwmu.Unlock()

	if mq.conn != nil && !mq.conn.IsClosed() {
		return mq.conn, nil
	}
	conn, err := amqp.Dial(mq.url)
	if err != nil {
		return nil, err
	}
	mq.conn = conn
	return mq.conn, nil
}
