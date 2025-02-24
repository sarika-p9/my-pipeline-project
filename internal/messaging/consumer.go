package messaging

import (
	"log"

	"github.com/streadway/amqp"
)

func StartConsumer(ch *amqp.Channel, queueName string) error {
	msgs, err := ch.Consume(
		queueName, "", true, false, false, false, nil,
	)
	if err != nil {
		return err
	}

	go func() {
		for msg := range msgs {
			log.Printf("📥 Consumed message: %s", msg.Body)
			// TODO: Process the message (e.g., trigger pipeline actions)
		}
	}()
	log.Println("✅ Consumer is running...")
	return nil
}
