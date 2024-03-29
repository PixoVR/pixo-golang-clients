package heartbeat

import (
	"encoding/json"
	"errors"
	abstract_client "github.com/PixoVR/pixo-golang-clients/pixo-platform/abstract-client"
)

type Pulse struct {
	SessionID int `json:"sessionId"`
}

func (c *client) SendPulse(sessionID int) error {
	body, err := json.Marshal(Pulse{SessionID: sessionID})
	if err != nil {
		return err
	}

	path := "pulse"

	httpRes, err := c.Post(path, body)
	if err != nil {
		return err
	}

	var res abstract_client.Response
	if err = json.Unmarshal(httpRes.Body(), &res); err != nil {
		return err
	}

	if httpRes.IsError() {
		return errors.New(res.Error)
	}

	return nil
}
