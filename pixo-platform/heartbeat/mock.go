package heartbeat

import (
	"context"
	"github.com/PixoVR/pixo-golang-clients/pixo-platform/abstract"
)

var _ Client = &MockClient{}

type MockClient struct {
	abstract.MockAbstractClient
	NumCalledPulse int
	PulseError     error
}

// SendPulse is a method that increments the number of times it was called and returns an error if PulseError is not nil
func (m *MockClient) SendPulse(ctx context.Context, sessionID int) error {
	m.NumCalledPulse++
	return m.PulseError
}

// SendPulsesWithCancel is a method that returns a channel and a cancel function
func (m *MockClient) SendPulsesWithCancel(ctx context.Context, sessionID int, periodSeconds float64) (chan error, context.CancelFunc) {
	errCh := make(chan error, 5)

	if m.PulseError != nil {
		errCh <- m.SendPulse(ctx, sessionID)
	}

	return errCh, func() {
		// Do nothing
	}
}
