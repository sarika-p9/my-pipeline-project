package messaging

import (
	"log"

	"github.com/streadway/amqp"
)

type Producer struct {
	channel *amqp.Channel
	queue   amqp.Queue
}

func NewProducer(ch *amqp.Channel, queueName string) (*Producer, error) {
	queue, err := ch.QueueDeclare(
		queueName, true, false, false, false, nil,
	)
	if err != nil {
		return nil, err
	}

	return &Producer{channel: ch, queue: queue}, nil
}

func (p *Producer) Publish(message string) error {
	err := p.channel.Publish(
		"", p.queue.Name, false, false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(message),
		},
	)
	if err != nil {
		return err
	}
	log.Printf("📤 Published message: %s", message)
	return nil
}
