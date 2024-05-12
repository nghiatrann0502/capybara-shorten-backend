package model

import (
	"encoding/json"
	"errors"
	"strings"
	"time"
)

type TrackingPayload struct {
	Name string        `json:"name"`
	Data *TrackingData `json:"data"`
}

type TrackingData struct {
	Id        int    `json:"id"`
	UserAgent string `json:"user_agent"`
	Referer   string `json:"referer"`
}

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

type URLCached struct {
	Id  int    `json:"id" redis:"id"`
	Url string `json:"url" redis:"url"`
}

func (u *URLCached) MarshalBinary() ([]byte, error) {
	return json.Marshal(u)
}

func (data *CreateShorten) Validate() error {
	if strings.TrimSpace(data.Url) == "" {
		return errors.New("url cannot be blank")
	}

	if strings.Split(data.Url, "://")[0] != "https" && strings.Split(data.Url, "://")[0] != "http" {
		return errors.New("url must start with http or https")
	}

	return nil
}
