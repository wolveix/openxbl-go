package openxbl

import (
	"errors"
	"time"
)

type Alert struct {
	ID            string    `json:"id"`
	Action        string    `json:"action"`
	Path          string    `json:"path"`
	ActorXuid     string    `json:"actorXuid"`
	ActorGamertag string    `json:"actorGamertag"`
	ParentType    string    `json:"parentType"`
	ParentPath    string    `json:"parentPath"`
	OwnerXuid     string    `json:"ownerXuid"`
	OwnerGamertag string    `json:"ownerGamertag"`
	Timestamp     time.Time `json:"timestamp"`
	Seen          bool      `json:"seen"`
	RootPath      string    `json:"rootPath"`
	ClubID        string    `json:"clubId"`
}

func (c *Client) GetAlerts() ([]*Alert, error) {
	response := struct {
		Alerts []*Alert `json:"alerts"`
	}{}

	if _, err := c.makeRequest("GET", "alerts", nil, &response); err != nil {
		return nil, err
	}

	if len(response.Alerts) == 0 {
		return nil, errors.New("failed to find alerts")
	}

	return response.Alerts, nil
}
