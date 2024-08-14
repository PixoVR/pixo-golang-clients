package headset

import (
	"context"
	abstract "github.com/PixoVR/pixo-golang-clients/pixo-platform/abstract-client"
)

type MockClient struct {
	abstract.MockAbstractClient

	NumCalledStartSession int
	StartSessionError     error

	NumCalledSendEvent int
	SendEventError     error

	NumCalledEndSession int
	EndSessionError     error
}

func (m *MockClient) StartSession(ctx context.Context, request EventRequest) (*EventResponse, error) {
	m.NumCalledStartSession++
	return nil, m.StartSessionError
}

func (m *MockClient) SendEvent(ctx context.Context, request EventRequest) (*EventResponse, error) {
	m.NumCalledSendEvent++
	return nil, m.SendEventError
}

func (m *MockClient) EndSession(ctx context.Context, request EventRequest) (*EventResponse, error) {
	m.NumCalledEndSession++
	return nil, m.EndSessionError
}
