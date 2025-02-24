package services

import (
	"log"

	"github.com/nats-io/nats.go"
)

func SubscribeToEvents(nc *nats.Conn, subject string) {
	_, err := nc.Subscribe(subject, func(m *nats.Msg) {
		log.Printf("📥 Real-time event received on '%s': %s", subject, string(m.Data))
	})
	if err != nil {
		log.Fatalf("❌ Failed to subscribe to events: %v", err)
	}
	log.Printf("✅ Subscribed to '%s' for real-time updates", subject)
}
