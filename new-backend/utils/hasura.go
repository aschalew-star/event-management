package utils

import (
	"bytes"
	"encoding/json"
	"net/http"
	"os"
)

type HasuraEvent struct {
	Name string      `json:"name"`
	Data interface{} `json:"data"`
}

func TriggerHasuraEvent(eventName string, data interface{}) {
	hasuraEndpoint := os.Getenv("HASURA_ENDPOINT")
	adminSecret := os.Getenv("HASURA_ADMIN_SECRET")

	if hasuraEndpoint == "" {
		return
	}

	event := HasuraEvent{
		Name: eventName,
		Data: data,
	}

	jsonData, err := json.Marshal(event)
	if err != nil {
		return
	}

	req, err := http.NewRequest("POST", hasuraEndpoint+"/event", bytes.NewBuffer(jsonData))
	if err != nil {
		return
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-hasura-admin-secret", adminSecret)

	client := &http.Client{}
	_, _ = client.Do(req)
}
