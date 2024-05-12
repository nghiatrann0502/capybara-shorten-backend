package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/nghiatrann0502/capybara-shorten-backend/internal/url-shorten/model"
	"github.com/nghiatrann0502/capybara-shorten-backend/internal/url-shorten/stores"
	"github.com/redis/go-redis/v9"
	"github.com/rs/cors"
	"io"
	"log"
	"net/http"
)

type Handler struct {
	store  stores.Store
	engine *http.ServeMux
}

type SimpleResponse struct {
	Success bool        `json:"success"`
	Error   string      `json:"error,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

func New(store stores.Store) (*Handler, error) {
	mux := http.NewServeMux()

	h := &Handler{
		store:  store,
		engine: mux,
	}

	return h, nil
}

func (h *Handler) InitCors() (http.Handler, error) {
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	})

	handler := c.Handler(h.engine)
	return handler, nil
}

func (h *Handler) Listen(handler http.Handler) error {
	if err := h.setHandlers(); err != nil {
		return err
	}

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", "8000"),
		Handler: handler,
	}

	err := srv.ListenAndServe()
	return err

}

func (h *Handler) CloseStore() error {
	return h.store.Storage.Close()
}

func (h *Handler) SimpleResponse(w http.ResponseWriter, status int, data interface{}, headers ...http.Header) error {
	var response SimpleResponse
	response.Success = true
	response.Data = data

	output, err := json.Marshal(response)
	if err != nil {
		return err
	}

	if len(headers) > 0 {
		for k, v := range headers[0] {
			w.Header()[k] = v
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	_, err = w.Write(output)
	if err != nil {
		return err
	}

	return nil
}

func (h *Handler) SimpleError(w http.ResponseWriter, status int, message string, headers ...http.Header) error {
	var response SimpleResponse
	response.Success = true
	response.Error = message

	output, err := json.Marshal(response)
	if err != nil {
		return err
	}

	if len(headers) > 0 {
		for k, v := range headers[0] {
			w.Header()[k] = v
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	_, err = w.Write(output)
	if err != nil {
		return err
	}

	return nil
}

func (h *Handler) Redirect(w http.ResponseWriter, r *http.Request, url string) {
	http.Redirect(w, r, url, http.StatusSeeOther)
}

func (h *Handler) ReadJSON(w http.ResponseWriter, r *http.Request, data interface{}) error {
	maxBytes := 1048576 // One megabyte
	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))

	dec := json.NewDecoder(r.Body)
	err := dec.Decode(data)
	if err != nil {
		return err
	}

	err = dec.Decode(&struct{}{})
	if err != io.EOF {
		return errors.New("body must have only a singble JSON value")
	}

	return nil
}

func (h *Handler) setHandlers() error {
	h.engine.HandleFunc("GET /ping", func(w http.ResponseWriter, r *http.Request) {
		h.SimpleResponse(w, http.StatusOK, "pong")
		return
	})

	h.engine.HandleFunc("POST /api/v1/url-shorten", h.createShortenHandler)
	h.engine.HandleFunc("GET /api/v1/url-shorten/{shortId}", h.getLongUrlHandler)
	return nil
}

func (h *Handler) createShortenHandler(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Url string `json:"url"`
	}

	if err := h.ReadJSON(w, r, &body); err != nil {
		h.SimpleError(w, http.StatusBadRequest, "invalid request")
		return
	}

	var shortId = h.store.GenerateShortLink(body.Url)

	data := model.CreateShorten{
		Url:     body.Url,
		ShortId: shortId,
	}
	if err := data.Validate(); err != nil {
		h.SimpleError(w, http.StatusBadRequest, err.Error())
		return
	}
	_, err := h.store.Storage.CreateShortURL(data)
	if err != nil {
		h.SimpleError(w, http.StatusInternalServerError, err.Error())
		return
	}

	h.SimpleResponse(w, http.StatusOK, shortId)
}

func (h *Handler) getLongUrlHandler(w http.ResponseWriter, r *http.Request) {
	shortId := r.PathValue("shortId")
	result := model.URLCached{}

	err := h.store.RedisClient.HMGet(context.TODO(), fmt.Sprintf("URL:%s", shortId), "id", "url").Scan(&result)
	if errors.Is(err, redis.Nil) {
		data, err := h.store.Storage.GetLongUrl(shortId)
		if err != nil {
			h.SimpleError(w, http.StatusInternalServerError, err.Error())
			return
		}

		if err := h.store.RedisClient.HSet(context.Background(), fmt.Sprintf("URL:%s", shortId), model.URLCached{Id: data.Id, Url: data.Url}).Err(); err != nil {
			log.Println(err)
		}

		go func() {
			err := h.store.PushToQueue("INCREASE_COUNT", &model.TrackingData{
				Id:        data.Id,
				Referer:   "localhost",
				UserAgent: "user_agent",
			})
			if err != nil {
				log.Println(err)
				fmt.Println("Error pushing to queue")
			}
			fmt.Println("done")
		}()
		h.Redirect(w, r, data.Url)
		return
	} else if err != nil {
		log.Println(err, "redis error")
		h.SimpleError(w, http.StatusInternalServerError, err.Error())
		return
	} else {
		go func() {
			err := h.store.PushToQueue("INCREASE_COUNT", &model.TrackingData{
				Id:        result.Id,
				Referer:   "localhost",
				UserAgent: "user_agent",
			})
			if err != nil {
				log.Println(err)
				fmt.Println("Error pushing to queue")
			}
		}()
		h.Redirect(w, r, result.Url)
	}
}
