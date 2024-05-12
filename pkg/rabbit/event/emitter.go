package event

import (
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
)

type Emitter struct {
	connection *amqp.Connection
}

func (e *Emitter) setup() error {
	channel, err := e.connection.Channel()
	if err != nil {
		log.Println(err, "err")
		return err
	}

	defer channel.Close()
	return DeclareExchange(channel)
}

func (e *Emitter) Push(event string, severity string) error {
	channel, err := e.connection.Channel()
	if err != nil {
		return err
	}

	defer channel.Close()

	log.Println("Pushing event", event, "with severity", severity)

	if err := channel.Publish(
		"logs_topic",
		severity,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(event),
		},
	); err != nil {
		return err
	}
	return nil
}

func NewEvenEmitter(conn *amqp.Connection) (Emitter, error) {
	e := Emitter{
		connection: conn,
	}

	if err := e.setup(); err != nil {
		return Emitter{}, err
	}
	return e, nil
}
