package platform

import (
	"context"
	"encoding/json"
	"time"
)

type Event struct {
	ID        int         `json:"id,omitempty"`
	SessionID int         `json:"sessionId,omitempty"`
	Type      string      `json:"eventType,omitempty"`
	Payload   interface{} `json:"jsonData,omitempty"`

	CreatedAt time.Time `json:"createdAt,omitempty"`
	UpdatedAt time.Time `json:"updatedAt,omitempty"`
}

type CreateEventResponse struct {
	Event Event `json:"createEvent"`
}

type UpdateEventResponse struct {
	Event Event `json:"updateEvent"`
}

type EventResponse struct {
	Event Event `json:"event"`
}

func (g *PlatformClient) GetEvent(ctx context.Context, id int) (*Event, error) {
	query := `query event($id: ID!) { event(id: $id) { id session } }`

	variables := map[string]interface{}{
		"id": id,
	}

	res, err := g.Client.ExecRaw(ctx, query, variables)
	if err != nil {
		return nil, err
	}

	var sessionResponse EventResponse
	if err = json.Unmarshal(res, &sessionResponse); err != nil {
		return nil, err
	}

	return &sessionResponse.Event, nil
}

func (g *PlatformClient) CreateEvent(ctx context.Context, event Event) (*Event, error) {
	query := `mutation createEvent($input: EventInput!) { createEvent(input: $input) { id } }`

	variables := map[string]interface{}{
		"input": map[string]interface{}{
			"sessionId": event.SessionID,
			"eventType": event.Type,
		},
	}

	if event.Payload != nil {
		variables["input"].(map[string]interface{})["jsonData"] = event.Payload
	} else {
		variables["input"].(map[string]interface{})["jsonData"] = "{}"
	}

	res, err := g.Client.ExecRaw(ctx, query, variables)
	if err != nil {
		return nil, err
	}

	var eventResponse CreateEventResponse
	if err = json.Unmarshal(res, &eventResponse); err != nil {
		return nil, err
	}

	return &eventResponse.Event, nil
}

func (g *PlatformClient) UpdateEvent(ctx context.Context, session Event) (*Event, error) {
	query := `mutation updateEvent($input: EventInput!) { updateEvent(input: $input) { id sessionId } }`

	variables := map[string]interface{}{
		"input": map[string]interface{}{
			"id": session.ID,
		},
	}

	res, err := g.Client.ExecRaw(ctx, query, variables)
	if err != nil {
		return nil, err
	}

	var sessionResponse UpdateEventResponse
	if err = json.Unmarshal(res, &sessionResponse); err != nil {
		return nil, err
	}

	return &sessionResponse.Event, nil
}