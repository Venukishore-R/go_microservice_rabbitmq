package events

import "github.com/rabbitmq/amqp091-go"

func declareExchange(ch *amqp091.Channel) error {
	return ch.ExchangeDeclare(
		"logs_topic", //name
		"topic",      //kind
		true,         //durable
		false,        //auto-deleted
		false,        //internal
		false,        //false
		nil,          //args
	)
}

func declareRandomQueue(ch *amqp091.Channel) (amqp091.Queue, error) {
	return ch.QueueDeclare(
		"microservice", //name
		false,          //durable
		false,          //auto-deleted
		true,           //exclusive
		false,          //nowait
		nil,            //args
	)
}
