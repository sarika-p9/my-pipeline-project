package messaging

import (
	"log"

	"github.com/streadway/amqp"
)

func ConnectRabbitMQ(url string) (*amqp.Connection, *amqp.Channel, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, nil, err
	}

	// ✅ Log connection success ONLY here
	log.Println("✅ Connected to RabbitMQ")
	return conn, ch, nil
}
