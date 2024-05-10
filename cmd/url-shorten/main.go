package main

import (
	"errors"
	"fmt"
	"github.com/nghiatrann0502/capybara-shorten-backend/internal/url-shorten/handlers"
	"github.com/nghiatrann0502/capybara-shorten-backend/internal/url-shorten/stores"
	"log"
	"os"
	"os/signal"
)

func main() {
	fmt.Println("Hello, World!")
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	closeFn, err := initURLShorten()
	if err != nil {
		log.Panicf("could not init shortener: %v", err)
	}
	<-stop
	log.Println("Shutting down...")
	closeFn()
}

func initURLShorten() (func(), error) {
	fmt.Println("Init URL Shorten")
	store, err := stores.New()
	if err != nil {
		log.Println(err.Error())
		return nil, errors.New("could not create store")
	}

	handler, err := handlers.New(*store)
	if err != nil {
		return nil, errors.New("could not create handlers")
	}

	server, err := handler.InitCors()
	if err != nil {
		return nil, errors.New("could not init cors")
	}

	go func() {
		if err := handler.Listen(server); err != nil {
			log.Panicf("could not listen to http handlers: %v", err)
		}
	}()
	return func() {
		if err = handler.CloseStore(); err != nil {
			log.Panicf("failed to stop the handlers: %v", err)
		}
	}, nil
}
