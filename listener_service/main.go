package main

import (
	"listener_service/events"
	"log"
	"math"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	// try to connect to rabbitmq server
	rabbitConn, err := connect()
	if err != nil {
		log.Fatalf("error while connection to rabbitmq: %v", err)
		return
	}

	defer rabbitConn.Close()

	// start listening to messages
	log.Println("listening to rabbitmq messages")

	// create consumer
	consumer, err := events.NewConsumer(rabbitConn)
	if err != nil {
		log.Fatalf("error while creating consumer: %v", err)
		return
	}

	// watch the queue and consume events
	if err = consumer.Listen([]string{"log.INFO", "log.ERROR", "log.WARNING"}); err != nil {
		log.Fatalf("error while listening: %v", err)
		return
	}
}

func connect() (*amqp.Connection, error) {
	var counts int64
	var backoff = 1 * time.Second
	var connection *amqp.Connection

	for {
		c, err := amqp.Dial("amqp://guest:guest@rabbitmq")
		if err != nil {
			log.Println("rabbitmq not yet ready-count: ", counts)
			counts++
		} else {
			log.Println("connected to rabbitmq server")
			connection = c
			break
		}

		if counts > 5 {
			log.Println("unable to connect to rabbitmq: ", err)
			return nil, err
		}

		backoff = time.Duration(math.Pow(float64(counts), 2)) * time.Second
		log.Println("backing off.....")
		time.Sleep(backoff)
		continue
	}

	return connection, nil
}
