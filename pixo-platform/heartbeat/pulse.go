package heartbeat

import (
	"encoding/json"
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

	_, err = c.Post(path, body)
	return err
}
