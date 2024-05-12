package event

import (
	"encoding/json"
	"fmt"
	"github.com/nghiatrann0502/capybara-shorten-backend/internal/logger/model"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
)

type Consumer struct {
	conn      *amqp.Connection
	queueName string
}

func NewConsumer(conn *amqp.Connection) (Consumer, error) {
	consumer := Consumer{
		conn: conn,
	}

	err := consumer.setup()
	if err != nil {
		return Consumer{}, err
	}

	return consumer, nil
}

func (consumer *Consumer) setup() error {
	channel, err := consumer.conn.Channel()
	if err != nil {
		return err
	}

	return DeclareExchange(channel)
}

func (consumer *Consumer) Listen(topic []string) error {
	ch, err := consumer.conn.Channel()
	if err != nil {
		return err
	}

	defer ch.Close()

	q, err := DeclareRandomQueue(ch)
	if err != nil {
		return err
	}

	for _, t := range topic {
		if err := ch.QueueBind(
			q.Name,
			t,
			"logs_topic",
			false,
			nil,
		); err != nil {
			return err
		}
	}

	messages, err := ch.Consume(q.Name, "", true, false, false, false, nil)
	if err != nil {
		return err
	}

	forever := make(chan bool)
	go func() {
		for d := range messages {
			var payload model.TrackingPayload
			_ = json.Unmarshal(d.Body, &payload)
			// Do something with the payload
			go handlePayload(consumer, payload)
		}
	}()

	fmt.Printf("Waiting for messages [Exchange, Queue] [logs_topic, %s]\n", q.Name)
	<-forever
	return nil
}

func handlePayload(c *Consumer, payload model.TrackingPayload) {
	switch payload.Name {
	case "INCREASE_COUNT":
		if err := handleCreateLog(c, payload.Data); err != nil {
			fmt.Println(err)
		}

	default:
		fmt.Println("Unknown event")
	}
}

func handleCreateLog(c *Consumer, payload model.TrackingData) error {
	data := &model.LogEntry{
		UrlId:     payload.Id,
		Referer:   payload.Referer,
		UserAgent: payload.UserAgent,
	}
	log.Println(data)
	//
	//if err := Storage.Insert(data); err != nil {
	//	return err
	//}

	return nil
}
