package event

import (
	"encoding/json"
	"fmt"
	"github.com/nghiatrann0502/capybara-shorten-backend/internal/logger/model"
	amqp "github.com/rabbitmq/amqp091-go"
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

	return declareExchange(channel)
}

func (consumer *Consumer) Listen(topic []string) error {
	ch, err := consumer.conn.Channel()
	if err != nil {
		return err
	}

	defer ch.Close()

	q, err := declareRandomQueue(ch)
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
			go handlePayload(payload)
		}
	}()

	fmt.Printf("Waiting for messages [Exchange, Queue] [logs_topic, %s]\n", q.Name)
	<-forever
	return nil
}

func handlePayload(payload model.TrackingPayload) {
	switch payload.Name {
	case "INCREASE_COUNT":
		fmt.Println(payload.Data, "Payload")
	default:
		fmt.Println("Unknown event")
	}
}
