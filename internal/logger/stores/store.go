package stores

import (
	"errors"
	_ "github.com/go-sql-driver/mysql"
	"github.com/nghiatrann0502/capybara-shorten-backend/pkg/rabbit/event"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
	"os"
)

type Store struct {
	Rabbit *amqp.Connection
}

func New() (*Store, error) {

	// RabbitMQ connection
	rabbit, err := event.Connect()
	defer rabbit.Close()

	if err != nil {
		return nil, errors.New("could not connect to RabbitMQ")
	}

	// Create consumer
	consumer, err := event.NewConsumer(rabbit)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	// Watch the queue and consume messages
	if err := consumer.Listen([]string{"log.INFO"}); err != nil {
		log.Println(err)
	}

	return &Store{
		Rabbit: rabbit,
	}, nil
}

func (s *Store) CloseStore() error {
	if err := s.Rabbit.Close(); err != nil {
		return err
	}
	return nil
}
