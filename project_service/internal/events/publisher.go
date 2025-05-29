package events

import (
	"encoding/json"
	"project_service/internal/model"
	"time"

	"github.com/nats-io/nats.go"
)

type Publisher struct {
	nc *nats.Conn
}

func NewPublisher(nc *nats.Conn) *Publisher {
	return &Publisher{nc: nc}
}

func (p *Publisher) PublishProjectCreated(project *model.Project) error {
	event := map[string]interface{}{
		"event":     "projects.created",
		"timestamp": time.Now().UTC().Format(time.RFC3339),
		"data": map[string]interface{}{
			"project_id": project.ID,
			"client_id":  project.ClientID,
			"title":      project.Title,
			"status":     project.Status,
		},
	}

	data, err := json.Marshal(event)
	if err != nil {
		return err
	}

	return p.nc.Publish("projects.created", data)
}
