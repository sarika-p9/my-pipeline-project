package services

import (
	"log"

	"github.com/nats-io/nats.go"
)

func PublishRealTimeEvent(nc *nats.Conn, subject string, message string) {
	err := nc.Publish(subject, []byte(message))
	if err != nil {
		log.Printf("❌ Failed to publish real-time event: %v", err)
	} else {
		log.Printf("📡 Published real-time event on '%s': %s", subject, message)
	}
}
