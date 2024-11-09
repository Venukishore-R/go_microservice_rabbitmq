package main

import (
	"fmt"
	"log"
	"log/slog"
	"math"
	"net/http"
	"time"

	"github.com/rabbitmq/amqp091-go"
)

const webPort = "1700"

type Config struct {
	Rabbit *amqp091.Connection
}

func main() {
	rabbitConn, err := connect()
	if err != nil {
		log.Fatalf("error while connection to rabbitmq: %v", err)
		return
	}

	app := Config{
		Rabbit: rabbitConn,
	}

	slog.Info("Starting broker-service", "port", webPort)

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", webPort),
		Handler: app.routes(),
	}

	if err = srv.ListenAndServe(); err != nil {
		log.Fatalf("error while serving broker-service: %v", err)
	}
}

func connect() (*amqp091.Connection, error) {
	var counts int64
	var backoff = 1 * time.Second
	var connection *amqp091.Connection

	for {
		c, err := amqp091.Dial("amqp://guest:guest@rabbitmq")
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
