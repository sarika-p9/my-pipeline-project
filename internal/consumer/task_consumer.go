package consumer

import (
	"encoding/json"
	"log"

	"github.com/sarika-p9/my-pipeline-project/internal/infrastructure"
)

type Task struct {
	PipelineID string `json:"pipeline_id"`
	JobType    string `json:"job_type"`
	Payload    string `json:"payload"`
}

type Consumer struct {
	RabbitMQ *infrastructure.RabbitMQ
}

func NewConsumer(rabbitMQ *infrastructure.RabbitMQ) *Consumer {
	return &Consumer{RabbitMQ: rabbitMQ}
}

func (c *Consumer) StartConsuming() {
	msgs, err := c.RabbitMQ.Channel.Consume(
		c.RabbitMQ.Queue.Name,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("❌ Failed to register consumer: %v", err)
	}

	log.Println("🚀 Consumer started. Waiting for tasks...")

	for d := range msgs {
		var task Task
		if err := json.Unmarshal(d.Body, &task); err != nil {
			log.Printf("❌ Failed to unmarshal task: %v", err)
			d.Nack(false, false)
			continue
		}

		log.Printf("📥 Processing task: %+v", task)

		d.Ack(false)
	}
}
