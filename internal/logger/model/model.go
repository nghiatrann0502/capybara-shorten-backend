package model

type TrackingPayload struct {
	Name string       `json:"name"`
	Data TrackingData `json:"data"`
}

type TrackingData struct {
	Id        int    `json:"id"`
	UserAgent string `json:"user_agent"`
	Referer   string `json:"referer"`
}
