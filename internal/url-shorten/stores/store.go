package stores

import (
	"crypto/sha256"
	"errors"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/itchyny/base58-go"
	"github.com/nghiatrann0502/capybara-shorten-backend/internal/url-shorten/model"
	"github.com/nghiatrann0502/capybara-shorten-backend/internal/url-shorten/stores/mysql"
	redisClient "github.com/nghiatrann0502/capybara-shorten-backend/internal/url-shorten/stores/redis"
	"github.com/redis/go-redis/v9"
	"math/big"
	"os"
	"time"
)

type Storage interface {
	CreateShortURL(data model.CreateShorten) (int, error)
	GetLongUrl(shortId string) (string, error)
	Close() error
}

type Store struct {
	Storage     Storage
	RedisClient *redis.Client
	idLength    int
}

func New() (*Store, error) {
	mySqlStore, err := mysql.New()

	if err != nil {
		return nil, errors.New("could not create store")
	}

	client, err := redisClient.New()

	return &Store{
		Storage:     mySqlStore,
		RedisClient: client,
		idLength:    6,
	}, nil
}

func (s *Store) GenerateShortLink(url string) string {
	urlHashBytes := sha256Of(url + time.Now().String())
	generatedNumber := new(big.Int).SetBytes(urlHashBytes).Uint64()
	finalString := base58Encoded([]byte(fmt.Sprintf("%d", generatedNumber)))
	return finalString[:8]
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
