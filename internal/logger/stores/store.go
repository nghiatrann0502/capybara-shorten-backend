package stores

import (
	"encoding/json"
	"errors"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/nghiatrann0502/capybara-shorten-backend/internal/logger/model"
	"github.com/nghiatrann0502/capybara-shorten-backend/internal/logger/stores/mongo"
	"github.com/nghiatrann0502/capybara-shorten-backend/pkg/rabbit/event"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
)

type Storage interface {
	Insert(entry *model.LogEntry) error
}

type Store struct {
	Rabbit  *amqp.Connection
	Storage Storage
}

func New() (*Store, error) {

	// Connect to mongodb
	mongoStore, err := mongo.New()
	if err != nil {
		return nil, errors.New("could not create store")
	}

	// RabbitMQ connection
	rabbit, err := event.Connect()
	//defer rabbit.Close()

	if err != nil {
		return nil, errors.New("could not connect to RabbitMQ")
	}

	return &Store{
		Rabbit:  rabbit,
		Storage: mongoStore,
	}, nil
}

func (s *Store) CreateConsumer() error {
	//Create consumer
	//consumer, err := event.NewConsumer(s.Rabbit)
	//if err != nil {
	//	log.Println(err)
	//	os.Exit(1)
	//}

	//Watch the queue and consume messages
	if err := s.Listen([]string{"log.INFO"}); err != nil {
		return err
	}
	return nil
}

func (s *Store) Listen(topic []string) error {
	ch, err := s.Rabbit.Channel()
	if err != nil {
		return err
	}

	defer ch.Close()

	q, err := event.DeclareRandomQueue(ch)
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
			go handlePayload(s, payload)
		}
	}()

	fmt.Printf("Waiting for messages [Exchange, Queue] [logs_topic, %s]\n", q.Name)
	<-forever
	return nil
}

func handlePayload(s *Store, payload model.TrackingPayload) {
	switch payload.Name {
	case "INCREASE_COUNT":
		if err := handleCreateLog(s, payload.Data); err != nil {
			fmt.Println(err)
		}

	default:
		fmt.Println("Unknown event")
	}
}

func handleCreateLog(s *Store, payload model.TrackingData) error {
	data := &model.LogEntry{
		UrlId:     payload.Id,
		Referer:   payload.Referer,
		UserAgent: payload.UserAgent,
	}
	log.Println(data, "data")
	//
	if err := s.Storage.Insert(data); err != nil {
		return err
	}

	return nil
}

func (s *Store) CloseStore() error {
	if err := s.Rabbit.Close(); err != nil {
		return err
	}
	return nil
}
