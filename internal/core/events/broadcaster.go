package events

import (
	"encoding/json"
	"log"
)

var Broadcast = make(chan string, 100)

type StageStatusUpdate struct {
	StageID string `json:"stage_id"`
	Status  string `json:"status"`
}

func SendUpdate(stageID string, status string) {
	message, err := json.Marshal(StageStatusUpdate{StageID: stageID, Status: status})
	if err != nil {
		log.Println("Error marshalling JSON:", err)
		return
	}
	Broadcast <- string(message)
}
