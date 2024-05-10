package model

import (
	"time"
)

type URLShorten struct {
	ID        int        `json:"id"`
	ShortId   string     `json:"short_id"`
	Url       string     `json:"url"`
	CreatedAt *time.Time `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
}

type CreateShorten struct {
	ShortId string `json:"short_id"`
	Url     string `json:"url"`
}
