package headset

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/PixoVR/pixo-golang-clients/pixo-platform/platform"
	"github.com/rs/zerolog/log"
	"io"
)

// Response is the response object for requests to the headset API
type Response struct {
	Error   bool          `json:"error"`
	Message string        `json:"message"`
	Data    EventResponse `json:"data"`
}

// EventRequest is the request object for sending events to the platform
type EventRequest struct {
	UUID         *string                `json:"uuid,omitempty"`
	SessionID    *int                   `json:"sessionID,omitempty"`
	DeviceID     string                 `json:"deviceId,omitempty"`
	ModuleID     int                    `json:"moduleId,omitempty"`
	Type         string                 `json:"event_type,omitempty"`
	OtherType    string                 `json:"eventType,omitempty"`
	Payload      map[string]interface{} `json:"jsonData,omitempty"`
	OtherPayload string                 `json:"jsondata,omitempty"`
}

// EventResponse is the response object when sending events to the platform
type EventResponse struct {
	platform.Event
	LessonStatus *string `json:"lessonStatus,omitempty"`
}

// StartSession creates a PIXOVR_SESSION_JOINED event
func (c *client) StartSession(ctx context.Context, request EventRequest) (*EventResponse, error) {
	request.Type = "PIXOVR_SESSION_JOINED"
	return c.SendEvent(ctx, request)
}

// EndSession creates a PIXOVR_SESSION_COMPLETE event
func (c *client) EndSession(ctx context.Context, request EventRequest) (*EventResponse, error) {
	request.Type = "PIXOVR_SESSION_COMPLETE"
	return c.SendEvent(ctx, request)
}

// SendEvent sends an event to the platform
func (c *client) SendEvent(ctx context.Context, request EventRequest) (*EventResponse, error) {
	body, err := json.Marshal(request)
	if err != nil {
		log.Error().Err(err).Msg("Failed to marshal event request")
		return nil, err
	}

	path := "event"

	res, err := c.Post(context.TODO(), path, body)
	if err != nil {
		log.Error().Err(err).Msg("Failed to post event request")
		return nil, err
	}

	resBody, _ := io.ReadAll(res.Body)

	var response Response
	if err = json.Unmarshal(resBody, &response); err != nil {
		log.Error().Err(err).Msg("Failed to unmarshal create event response")
		return nil, err
	}

	if response.Error {
		log.Error().Msg(response.Message)
		return nil, errors.New(response.Message)
	}

	return &response.Data, nil
}
