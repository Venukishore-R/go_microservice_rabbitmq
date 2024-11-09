package events

import (
	"log"

	"github.com/rabbitmq/amqp091-go"
)

type Emitter struct {
	connection *amqp091.Connection
}

func (e *Emitter) setup() error {
	channel, err := e.connection.Channel()
	if err != nil {
		return err
	}

	defer channel.Close()

	return declareExchange(channel)
}

func (e *Emitter) Push(event string, severity string) error {
	channel, err := e.connection.Channel()
	if err != nil {
		return err
	}

	defer channel.Close()

	log.Println("publishing to channel")

	if err = channel.Publish(
		"logs_topic",
		severity,
		false,
		false,
		amqp091.Publishing{
			ContentType: "text/plain",
			Body:        []byte(event),
		},
	); err != nil {
		return err
	}

	return nil
}

func NewEventEmitter(connection *amqp091.Connection) (Emitter, error) {
	emitter := Emitter{
		connection: connection,
	}

	if err := emitter.setup(); err != nil {
		return Emitter{}, err
	}

	return emitter, nil
}
