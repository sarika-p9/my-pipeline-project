// internal/producer/task_producer.go

package producer

import (
	"encoding/json"
	"log"

	"github.com/sarika-p9/my-pipeline-project/internal/infrastructure"

	"github.com/streadway/amqp"
)

type Task struct {
	PipelineID string `json:"pipeline_id"`
	JobType    string `json:"job_type"`
	Payload    string `json:"payload"`
}

type Producer struct {
	RabbitMQ *infrastructure.RabbitMQ
}

func NewProducer(rabbitMQ *infrastructure.RabbitMQ) *Producer {
	return &Producer{RabbitMQ: rabbitMQ}
}

func (p *Producer) PublishTask(task Task) error {
	body, err := json.Marshal(task)
	if err != nil {
		return err
	}

	err = p.RabbitMQ.Channel.Publish(
		"",                    // Exchange
		p.RabbitMQ.Queue.Name, // Routing Key
		false,                 // Mandatory
		false,                 // Immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
	if err != nil {
		return err
	}

	log.Printf("✅ Published task to queue: %+v", task)
	return nil
}
