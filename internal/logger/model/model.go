package model

import (
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

var client *mongo.Client

type TrackingPayload struct {
	Name string       `json:"name"`
	Data TrackingData `json:"data"`
}

type TrackingData struct {
	Id        int    `json:"id"`
	UserAgent string `json:"user_agent"`
	Referer   string `json:"referer"`
}

type LogEntry struct {
	ID        string    `bson:"_id,omitempty" json:"id,omitempty"`
	UrlId     int       `bson:"url_id" json:"url_id"`
	Referer   string    `bson:"referer" json:"referer"`
	UserAgent string    `bson:"data" json:"data"`
	CreatedAt time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time `bson:"updated_at" json:"updated_at"`
}
