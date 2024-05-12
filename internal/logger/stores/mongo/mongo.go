package mongo

import (
	"context"
	"github.com/nghiatrann0502/capybara-shorten-backend/internal/logger/model"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"time"
)

type Store struct {
	db *mongo.Client
}

func New() (*Store, error) {
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	clientOptions.SetAuth(options.Credential{
		Username: "admin",
		Password: "password",
	})
	// Connect
	c, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Println("Error connecting to MongoDB", err)
		return nil, err
	}

	log.Println("Connected to mongo")

	//ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	//defer cancel()
	//
	//defer func() {
	//	if err := c.Disconnect(ctx); err != nil {
	//		log.Panic(err)
	//	}
	//}()

	return &Store{db: c}, nil
}

func (s *Store) Insert(entry *model.LogEntry) error {
	collection := s.db.Database("logs").Collection("logs")

	_, err := collection.InsertOne(context.TODO(),
		model.LogEntry{
			UserAgent: entry.UserAgent,
			Referer:   entry.Referer,
			UrlId:     entry.UrlId,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	)

	if err != nil {
		log.Println("Error inserting entry entry", err)
		return err
	}

	return nil
}
