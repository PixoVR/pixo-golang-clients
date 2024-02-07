package heartbeat

import (
	"encoding/json"
	"errors"
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

	if res, err := c.Post(path, body); err != nil {
		return errors.New(res.String())
	}

	return nil
}
