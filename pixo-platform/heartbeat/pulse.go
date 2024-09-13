package heartbeat

import (
	"context"
	"encoding/json"
	"errors"
	abstract_client "github.com/PixoVR/pixo-golang-clients/pixo-platform/abstract"
	"io"
	"net/http"
)

// Pulse is a struct that represents a pulse sent to the heartbeat service.
type Pulse struct {
	SessionID int `json:"sessionId"`
}

// SendPulse sends a pulse to the heartbeat service
func (c *client) SendPulse(sessionID int) error {
	body, err := json.Marshal(Pulse{SessionID: sessionID})
	if err != nil {
		return err
	}

	path := "pulse"

	httpRes, err := c.Post(context.TODO(), path, body)
	if err != nil {
		return err
	}

	resBody, _ := io.ReadAll(httpRes.Body)

	var res abstract_client.Response
	if err = json.Unmarshal(resBody, &res); err != nil {
		return err
	}

	if httpRes.StatusCode != http.StatusOK {
		return errors.New(res.Error)
	}

	return nil
}
