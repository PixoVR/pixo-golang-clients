package platform

import (
	"context"
	"encoding/json"
	"errors"
	"time"
)

type Event struct {
	ID          int      `json:"id,omitempty"`
	SessionID   *int     `json:"sessionId,omitempty"`
	SessionUUID *string  `json:"sessionUuid,omitempty"`
	Session     *Session `json:"session,omitempty"`
	Type        string   `json:"type,omitempty"`
	Payload     string   `json:"jsonData,omitempty"`

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

func (p *PlatformClient) GetEvent(ctx context.Context, id int) (*Event, error) {
	query := `query event($id: ID!) { event(id: $id) { id session } }`

	variables := map[string]interface{}{
		"id": id,
	}

	res, err := p.Client.ExecRaw(ctx, query, variables)
	if err != nil {
		return nil, err
	}

	var sessionResponse EventResponse
	if err = json.Unmarshal(res, &sessionResponse); err != nil {
		return nil, err
	}

	return &sessionResponse.Event, nil
}

func (p *PlatformClient) CreateEvent(ctx context.Context, event *Event) error {
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

	if event.Payload != "" {
		variables["input"].(map[string]interface{})["payload"] = event.Payload
	}

	res, err := p.Client.ExecRaw(ctx, query, variables)
	if err != nil {
		return err
	}

	var eventResponse CreateEventResponse
	if err = json.Unmarshal(res, &eventResponse); err != nil {
		return err
	}

	*event = eventResponse.Event
	return nil
}

func (p *PlatformClient) UpdateEvent(ctx context.Context, session Event) (*Event, error) {
	query := `mutation updateEvent($input: EventInput!) { updateEvent(input: $input) { id sessionId } }`

	variables := map[string]interface{}{
		"input": map[string]interface{}{
			"id": session.ID,
		},
	}

	res, err := p.Client.ExecRaw(ctx, query, variables)
	if err != nil {
		return nil, err
	}

	var sessionResponse UpdateEventResponse
	if err = json.Unmarshal(res, &sessionResponse); err != nil {
		return nil, err
	}

	return &sessionResponse.Event, nil
}
