package infrastructure

import (
	"log"

	"github.com/nats-io/nats.go"
)

func ConnectNATS(natsURL string) *nats.Conn {
	nc, err := nats.Connect(natsURL)
	if err != nil {
		log.Fatalf("❌ Failed to connect to NATS: %v", err)
	}
	log.Println("✅ Connected to NATS")
	return nc
}
