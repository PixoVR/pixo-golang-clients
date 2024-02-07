package heartbeat

import abstract_client "github.com/PixoVR/pixo-golang-clients/pixo-platform/abstract-client"

type MockClient struct {
	abstract_client.MockAbstractClient
	NumCalledPulse int
}

func (m *MockClient) SendPulse(sessionID int) error {
	m.NumCalledPulse++
	return nil
}
