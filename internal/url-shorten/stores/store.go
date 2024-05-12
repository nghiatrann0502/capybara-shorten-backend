package stores

import (
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/itchyny/base58-go"
	"github.com/nghiatrann0502/capybara-shorten-backend/internal/url-shorten/model"
	"github.com/nghiatrann0502/capybara-shorten-backend/internal/url-shorten/stores/mysql"
	redisClient "github.com/nghiatrann0502/capybara-shorten-backend/internal/url-shorten/stores/redis"
	"github.com/nghiatrann0502/capybara-shorten-backend/pkg/rabbit/event"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/redis/go-redis/v9"
	"math/big"
	"os"
	"time"
)

type Storage interface {
	CreateShortURL(data model.CreateShorten) (int, error)
	GetLongUrl(shortId string) (*model.URLCached, error)
	Close() error
}

type Store struct {
	Storage     Storage
	RedisClient *redis.Client
	Rabbit      *amqp.Connection
	idLength    int
}

func New() (*Store, error) {
	mySqlStore, err := mysql.New()

	if err != nil {
		return nil, errors.New("could not create store")
	}

	client, err := redisClient.New()

	if err != nil {
		return nil, errors.New("could not create redis client")
	}

	// RabbitMQ connection
	rabbit, err := event.Connect()
	//defer rabbit.Close()

	if err != nil {
		return nil, errors.New("could not connect to RabbitMQ")
	}

	return &Store{
		Storage:     mySqlStore,
		RedisClient: client,
		Rabbit:      rabbit,
		idLength:    6,
	}, nil
}

func (s *Store) GenerateShortLink(url string) string {
	urlHashBytes := sha256Of(url + time.Now().String())
	generatedNumber := new(big.Int).SetBytes(urlHashBytes).Uint64()
	finalString := base58Encoded([]byte(fmt.Sprintf("%d", generatedNumber)))
	return finalString[:8]
}

func (s *Store) PushToQueue(name string, data *model.TrackingData) error {

	e, err := event.NewEvenEmitter(s.Rabbit)
	if err != nil {
		return err
	}

	payload := &model.TrackingPayload{
		Name: name,
		Data: data,
	}

	j, _ := json.MarshalIndent(payload, "", "\t")
	if err := e.Push(string(j), "log.INFO"); err != nil {
		return err
	}

	return nil
}

func sha256Of(input string) []byte {
	algorithm := sha256.New()
	algorithm.Write([]byte(input))
	return algorithm.Sum(nil)
}

func base58Encoded(bytes []byte) string {
	encoding := base58.BitcoinEncoding
	encoded, err := encoding.Encode(bytes)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	return string(encoded)
}
