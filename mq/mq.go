package mq

import (
	"bytes"
	"context"
	"github.com/streadway/amqp"
	"log"
	"sync"
)

const (
	baseExchangeName = "clipOne"
)

type MsgManager struct {
	conn         *amqp.Connection
	ch           *amqp.Channel
	exchangeName string
	backData     []byte
	ReceiveCh    chan []byte
	url          string
	rwLock       sync.RWMutex
}

func NewMsgManager(user, url string) *MsgManager {
	return &MsgManager{
		exchangeName: baseExchangeName + user,
		ReceiveCh:    make(chan []byte, 10),
		url:          url,
	}
}

func (m *MsgManager) Init() error {
	conn, err := amqp.Dial(m.url)
	if err != nil {
		log.Println("Failed to connect to RabbitMQ")
		return err
	}
	m.conn = conn

	ch, err := conn.Channel()
	if err != nil {
		log.Println("Failed to open a channel")
		return err
	}
	m.ch = ch

	err = ch.ExchangeDeclare(
		m.exchangeName, // name
		"fanout",       // type
		true,           // durable
		false,          // auto-deleted
		false,          // internal
		false,          // no-wait
		nil,            // arguments
	)
	if err != nil {
		m.Close()
		log.Println("Failed to declare an exchange")
		return err
	}

	return nil
}

func (m *MsgManager) Close() {
	m.ch.Close()
	m.conn.Close()
}

func (m *MsgManager) Receive(ctx context.Context) error {
	q, err := m.ch.QueueDeclare(
		"",    // name
		false, // durable
		false, // delete when unused
		true,  // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		log.Println("Failed to declare a queue")
		return err
	}

	err = m.ch.QueueBind(
		q.Name,         // queue name
		"",             // routing key
		m.exchangeName, // exchange
		false,
		nil,
	)
	if err != nil {
		log.Println("Failed to bind a queue")
		return err
	}

	msgCh, err := m.ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)

	if err != nil {
		log.Println("Failed to register a consumer")
		return err
	}

	go func() {
		for {
			select {
			case d := <-msgCh:
				m.rwLock.RLock()
				if m.backData != nil && bytes.Equal(m.backData, d.Body) {
					m.rwLock.RUnlock()
					continue
				}
				m.rwLock.RUnlock()
				m.ReceiveCh <- d.Body
			case <-ctx.Done():
				log.Println(" [receive]: ctx done!")
				return
			}
		}
	}()
	return nil
}

func (m *MsgManager) Send(data []byte) error {
	err := m.ch.Publish(
		m.exchangeName, // exchange
		"",             // routing key
		false,          // mandatory
		false,          // immediate
		amqp.Publishing{
			Body: data,
		})
	if err != nil {
		log.Println("Failed to publish a message")
		return err
	}
	m.rwLock.Lock()
	m.backData = data
	m.rwLock.Unlock()
	return nil
}
