package heartbeat

import abstract_client "github.com/PixoVR/pixo-golang-clients/pixo-platform/abstract-client"

type MockClient struct {
	abstract_client.MockAbstractClient
	NumCalledPulse int
	PulseError     error
}

// SendPulse is a method that increments the number of times it was called and returns an error if PulseError is not nil
func (m *MockClient) SendPulse(sessionID int) error {
	m.NumCalledPulse++
	if m.PulseError != nil {
		return m.PulseError
	}
	return nil
}
