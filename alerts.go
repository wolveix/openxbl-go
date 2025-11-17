package openxbl

import (
	"context"
	"errors"
	"net/http"
	"time"
)

type Alert struct {
	ID            string    `json:"id"`
	Action        string    `json:"action"`
	Path          string    `json:"path"`
	ActorXUID     string    `json:"actorXuid"`
	ActorGamertag string    `json:"actorGamertag"`
	ParentType    string    `json:"parentType"`
	ParentPath    string    `json:"parentPath"`
	OwnerXUID     string    `json:"ownerXuid"`
	OwnerGamertag string    `json:"ownerGamertag"`
	Timestamp     time.Time `json:"timestamp"`
	Seen          bool      `json:"seen"`
	RootPath      string    `json:"rootPath"`
	ClubID        string    `json:"clubId"`
}

func (c *Client) GetAlerts(ctx context.Context) ([]*Alert, error) {
	response := struct {
		Alerts []*Alert `json:"alerts"`
	}{}

	if _, err := c.makeRequest(ctx, http.MethodGet, "alerts", nil, &response); err != nil {
		return nil, err
	}

	if len(response.Alerts) == 0 {
		return nil, errors.New("find alerts")
	}

	return response.Alerts, nil
}
