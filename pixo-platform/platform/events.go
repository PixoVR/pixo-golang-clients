package platform

import (
	"context"
	"encoding/json"
	"errors"
	"time"
)

type Event struct {
	ID          int                    `json:"id,omitempty"`
	SessionID   *int                   `json:"sessionId,omitempty"`
	SessionUUID *string                `json:"sessionUuid,omitempty"`
	Session     *Session               `json:"session,omitempty"`
	Type        string                 `json:"type,omitempty"`
	Payload     map[string]interface{} `json:"jsonData,omitempty"`

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

func (p *clientImpl) GetEvent(ctx context.Context, id int) (*Event, error) {
	query := `query event($id: ID!) { event(id: $id) { id session } }`

	variables := map[string]interface{}{
		"id": id,
	}

	var res EventResponse
	if err := p.Exec(ctx, query, &res, variables); err != nil {
		return nil, err
	}

	return &res.Event, nil
}

func (p *clientImpl) CreateEvent(ctx context.Context, event *Event) error {
	if event == nil {
		return errors.New("event is nil")
	}

	query := `mutation createEvent($input: EventInput!) { createEvent(input: $input) { id sessionId } }`

	variables := map[string]interface{}{
		"input": map[string]interface{}{
			"sessionId":   event.SessionID,
			"sessionUuid": event.SessionUUID,
			"type":        event.Type,
		},
	}

	if event.Payload != nil {
		payload, err := json.Marshal(event.Payload)
		if err != nil {
			return errors.New("invalid json")
		}

		variables["input"].(map[string]interface{})["payload"] = string(payload)
	}

	var res CreateEventResponse
	if err := p.Exec(ctx, query, &res, variables); err != nil {
		return err
	}

	*event = res.Event
	return nil
}

func (p *clientImpl) UpdateEvent(ctx context.Context, session Event) (*Event, error) {
	query := `mutation updateEvent($input: EventInput!) { updateEvent(input: $input) { id sessionId } }`

	variables := map[string]interface{}{
		"input": map[string]interface{}{
			"id": session.ID,
		},
	}

	var res UpdateEventResponse
	if err := p.Exec(ctx, query, &res, variables); err != nil {
		return nil, err
	}

	return &res.Event, nil
}
