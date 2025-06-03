package heartbeat

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/PixoVR/pixo-golang-clients/pixo-platform/abstract"
	"github.com/rs/zerolog/log"
	"io"
	"net/http"
	"time"
)

// Pulse is a struct that represents a pulse sent to the heartbeat service.
type Pulse struct {
	SessionID int `json:"sessionId"`
}

// SendPulse sends a pulse to the heartbeat service
func (c *client) SendPulse(ctx context.Context, sessionID int) error {
	body, err := json.Marshal(Pulse{SessionID: sessionID})
	if err != nil {
		return err
	}

	httpRes, err := c.Post(ctx, "pulse", body)
	if err != nil {
		return err
	}

	resBody, _ := io.ReadAll(httpRes.Body)

	var res abstract.Response
	if err = json.Unmarshal(resBody, &res); err != nil {
		return err
	}

	if httpRes.StatusCode != http.StatusOK {
		return errors.New(res.Error)
	}

	return nil
}

// SendPulsesWithCancel starts sending pulses in a new goroutine and returns a channel to receive errors and a cancel function
func (c *client) SendPulsesWithCancel(ctx context.Context, sessionID int, periodSeconds float64) (chan error, context.CancelFunc) {
	errCh := make(chan error, 5)
	newCtx, cancel := context.WithCancel(ctx)

	go c.startSending(newCtx, sessionID, periodSeconds, errCh)

	return errCh, cancel
}

// startSending starts sending pulses to the heartbeat service
func (c *client) startSending(ctx context.Context, sessionID int, periodSeconds float64, errCh chan error) {
	for {
		select {
		case <-ctx.Done():
			log.Debug().Msg("Stopping heartbeat for session")
			return
		default:
			if err := c.SendPulse(ctx, sessionID); err != nil {
				log.Error().Err(err).Msg("Error sending heartbeat")
				errCh <- err
			} else {
				log.Debug().Msg("Sent pulse to platform heartbeat service")
			}
			time.Sleep(time.Duration(periodSeconds) * time.Second)
		}
	}
}
