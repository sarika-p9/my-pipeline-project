package main

import (
	"log"

	"github.com/sarika-p9/my-pipeline-project/cmd/democtl/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		log.Fatalf("Error executing command: %v", err)
	}
}
