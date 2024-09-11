package headset

import (
	"context"
	abstract "github.com/PixoVR/pixo-golang-clients/pixo-platform/abstract-client"
	"github.com/PixoVR/pixo-golang-clients/pixo-platform/platform"
)

// MockClient is a mock implementation of the Client interface.
type MockClient struct {
	abstract.MockAbstractClient

	NumCalledStartSession int
	StartSessionError     error

	NumCalledSendEvent int
	SendEventError     error

	NumCalledEndSession int
	EndSessionError     error
}

// StartSession returns an error if provided, otherwise returns a mock EventResponse.
func (m *MockClient) StartSession(ctx context.Context, request EventRequest) (*EventResponse, error) {
	m.Lock.Lock()
	defer m.Lock.Unlock()

	m.NumCalledStartSession++

	if m.StartSessionError != nil {
		return nil, m.StartSessionError
	}

	return &EventResponse{
		Event: platform.Event{
			ID:        1,
			SessionID: &[]int{1}[0],
			Session: &platform.Session{
				ID: 1,
			},
			Type:    request.Type,
			Payload: request.Payload,
		},
	}, nil
}

// SendEvent returns an error if provided, otherwise returns a mock EventResponse.
func (m *MockClient) SendEvent(ctx context.Context, request EventRequest) (*EventResponse, error) {
	m.Lock.Lock()
	defer m.Lock.Unlock()

	m.NumCalledSendEvent++

	if m.SendEventError != nil {
		return nil, nil
	}

	return &EventResponse{
		Event: platform.Event{
			ID:        1,
			SessionID: &[]int{1}[0],
			Session: &platform.Session{
				ID: 1,
			},
			Type:    request.Type,
			Payload: request.Payload,
		},
	}, nil
}

// EndSession returns an error if provided, otherwise returns a mock EventResponse.
func (m *MockClient) EndSession(ctx context.Context, request EventRequest) (*EventResponse, error) {
	m.Lock.Lock()
	defer m.Lock.Unlock()

	m.NumCalledEndSession++

	if m.EndSessionError != nil {
		return nil, m.EndSessionError
	}

	return &EventResponse{
		Event: platform.Event{
			ID:        1,
			SessionID: &[]int{1}[0],
			Session: &platform.Session{
				ID: 1,
			},
			Type:    request.Type,
			Payload: request.Payload,
		},
	}, nil
}
