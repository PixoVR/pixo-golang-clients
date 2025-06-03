package headset

import (
	"context"
	abstract "github.com/PixoVR/pixo-golang-clients/pixo-platform/abstract"
	"github.com/PixoVR/pixo-golang-clients/pixo-platform/platform"
)

// MockClient is a mock implementation of the Client interface.
type MockClient struct {
	abstract.MockAbstractClient

	CalledStartSessionWith []EventRequest
	StartSessionError      error

	CalledSendEventWith []EventRequest
	SendEventError      error

	CalledEndSessionWith []EventRequest
	EndSessionError      error
}

// StartSession returns an error if provided, otherwise returns a mock EventResponse.
func (m *MockClient) StartSession(ctx context.Context, request EventRequest) (*EventResponse, error) {
	m.Lock.Lock()
	defer m.Lock.Unlock()

	m.CalledStartSessionWith = append(m.CalledStartSessionWith, request)

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
			Type:       request.Type,
			PayloadMap: request.Payload,
		},
	}, nil
}

// SendEvent returns an error if provided, otherwise returns a mock EventResponse.
func (m *MockClient) SendEvent(ctx context.Context, request EventRequest) (*EventResponse, error) {
	m.Lock.Lock()
	defer m.Lock.Unlock()

	m.CalledSendEventWith = append(m.CalledSendEventWith, request)

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
			Type:       request.Type,
			PayloadMap: request.Payload,
		},
	}, nil
}

// EndSession returns an error if provided, otherwise returns a mock EventResponse.
func (m *MockClient) EndSession(ctx context.Context, request EventRequest) (*EventResponse, error) {
	m.Lock.Lock()
	defer m.Lock.Unlock()

	m.CalledEndSessionWith = append(m.CalledEndSessionWith, request)

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
			Type:       request.Type,
			PayloadMap: request.Payload,
		},
	}, nil
}
