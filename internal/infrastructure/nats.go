package infrastructure

import (
	"log"
	"time"

	"github.com/nats-io/nats.go"
)

func ConnectNATS(natsURL string) *nats.Conn {
	var nc *nats.Conn
	var err error

	maxRetries := 5
	for i := 1; i <= maxRetries; i++ {
		nc, err = nats.Connect(natsURL)
		if err == nil {
			log.Println("✅ Connected to NATS")
			return nc
		}

		log.Printf("❌ Failed to connect to NATS (attempt %d/%d): %v", i, maxRetries, err)
		time.Sleep(time.Duration(i) * time.Second) // Exponential backoff
	}

	log.Fatalf("❌ Could not connect to NATS after %d attempts: %v", maxRetries, err)
	return nil
}
