package headset

import (
	"context"
	abstract "github.com/PixoVR/pixo-golang-clients/pixo-platform/abstract-client"
	"github.com/PixoVR/pixo-golang-clients/pixo-platform/platform"
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

func (m *MockClient) SendEvent(ctx context.Context, request EventRequest) (*EventResponse, error) {
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

func (m *MockClient) EndSession(ctx context.Context, request EventRequest) (*EventResponse, error) {
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
