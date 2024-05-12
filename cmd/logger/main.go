package main

import (
	"errors"
	"fmt"
	"github.com/nghiatrann0502/capybara-shorten-backend/internal/logger/stores"
	"log"
	"os"
	"os/signal"
)

func main() {
	fmt.Println("Hello, World!")
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	closeFn, err := initLogger()
	if err != nil {
		log.Panicf("could not init shortener: %v", err)
	}
	<-stop
	log.Println("Shutting down...")
	closeFn()
}

func initLogger() (func(), error) {
	fmt.Println("Init Logger")
	store, err := stores.New()
	if err != nil {
		log.Println(err.Error())
		return nil, errors.New("could not create store")
	}
	if err := store.CreateConsumer(); err != nil {
		log.Println(err.Error())
		return nil, errors.New("could not create consumer")
	}

	return func() {
		if err = store.CloseStore(); err != nil {
			log.Panicf("failed to stop the handlers: %v", err)
		}
	}, nil
}
