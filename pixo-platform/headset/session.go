package headset

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/PixoVR/pixo-golang-clients/pixo-platform/platform"
	"github.com/rs/zerolog/log"
)

type Response struct {
	Error   bool          `json:"error"`
	Message string        `json:"message"`
	Data    EventResponse `json:"data"`
}

type EventRequest struct {
	UUID         *string                `json:"uuid,omitempty"`
	SessionID    *int                   `json:"sessionID,omitempty"`
	DeviceID     string                 `json:"deviceId,omitempty"`
	ModuleID     int                    `json:"moduleId,omitempty"`
	Type         string                 `json:"event_type,omitempty"`
	OtherType    string                 `json:"eventType,omitempty"`
	Payload      map[string]interface{} `json:"jsonData,omitempty"`
	OtherPayload map[string]interface{} `json:"jsondata,omitempty"`
}

type EventResponse struct {
	platform.Event
	LessonStatus *string `json:"lessonStatus,omitempty"`
}

func (c *client) StartSession(ctx context.Context, request EventRequest) (*EventResponse, error) {
	request.Type = "PIXOVR_SESSION_JOINED"

	body, err := json.Marshal(request)
	if err != nil {
		log.Error().Err(err).Msg("Failed to marshal allocate server request")
		return nil, err
	}

	path := "event"

	res, err := c.Post(path, body)
	if err != nil {
		log.Error().Err(err).Msg("Failed to post allocate server request")
		return nil, err
	}

	var response Response
	if err = json.Unmarshal(res.Body(), &response); err != nil {
		log.Error().Err(err).Msg("Failed to unmarshal allocate server response")
		return nil, err
	}

	if response.Error {
		log.Error().Msg(response.Message)
		return nil, errors.New(response.Message)
	}

	return &response.Data, nil
}

func (c *client) EndSession(ctx context.Context, request EventRequest) (*EventResponse, error) {
	request.Type = "PIXOVR_SESSION_COMPLETE"

	body, err := json.Marshal(request)
	if err != nil {
		log.Error().Err(err).Msg("Failed to marshal allocate server request")
		return nil, err
	}

	path := "event"

	res, err := c.Post(path, body)
	if err != nil {
		log.Error().Err(err).Msg("Failed to post allocate server request")
		return nil, err
	}

	var response Response
	if err = json.Unmarshal(res.Body(), &response); err != nil {
		log.Error().Err(err).Msg("Failed to unmarshal allocate server response")
		return nil, err
	}

	if response.Error {
		log.Error().Msg(response.Message)
		return nil, errors.New(response.Message)
	}

	return &response.Data, nil
}
