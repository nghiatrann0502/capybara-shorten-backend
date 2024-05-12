package model

import (
	"errors"
	"strings"
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

func (data *CreateShorten) Validate() error {
	if strings.TrimSpace(data.Url) == "" {
		return errors.New("url cannot be blank")
	}

	if strings.Split(data.Url, "://")[0] != "https" && strings.Split(data.Url, "://")[0] != "http" {
		return errors.New("url must start with http or https")
	}

	return nil
}
