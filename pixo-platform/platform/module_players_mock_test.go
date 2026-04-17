package platform_test

import (
	"context"
	"errors"
	"testing"

	"github.com/PixoVR/pixo-golang-clients/pixo-platform/platform"
)

func TestMockClient_GetModulePlayers_Default(t *testing.T) {
	mock := &platform.MockClient{}

	modulePlayers, err := mock.GetModulePlayers(context.Background())

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(modulePlayers) != 1 {
		t.Fatalf("expected 1 module player, got %d", len(modulePlayers))
	}
	got := modulePlayers[0]
	if got.ID != 1 {
		t.Errorf("expected ID 1, got %d", got.ID)
	}
	if got.Name != "test-module-player" {
		t.Errorf("expected Name test-module-player, got %q", got.Name)
	}
	if got.DistributorID != 1 {
		t.Errorf("expected DistributorID 1, got %d", got.DistributorID)
	}
	if got.LaunchProtocol != "pixo" {
		t.Errorf("expected LaunchProtocol pixo, got %q", got.LaunchProtocol)
	}
	if got.Distributor.ID != 1 || got.Distributor.Name != "test-org" {
		t.Errorf("unexpected distributor: %+v", got.Distributor)
	}
	if mock.NumCalledGetModulePlayers != 1 {
		t.Errorf("expected NumCalledGetModulePlayers 1, got %d", mock.NumCalledGetModulePlayers)
	}
}

func TestMockClient_GetModulePlayers_ReturnsConfigured(t *testing.T) {
	expected := []platform.ModulePlayer{
		{ID: 42, Name: "custom", LaunchProtocol: "webxr"},
	}
	mock := &platform.MockClient{GetModulePlayersReturn: expected}

	modulePlayers, err := mock.GetModulePlayers(context.Background())

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(modulePlayers) != 1 || modulePlayers[0].ID != 42 || modulePlayers[0].Name != "custom" || modulePlayers[0].LaunchProtocol != "webxr" {
		t.Errorf("unexpected result: %+v", modulePlayers)
	}
	if mock.NumCalledGetModulePlayers != 1 {
		t.Errorf("expected NumCalledGetModulePlayers 1, got %d", mock.NumCalledGetModulePlayers)
	}
}

func TestMockClient_GetModulePlayers_ReturnsError(t *testing.T) {
	mockErr := errors.New("boom")
	mock := &platform.MockClient{GetModulePlayersError: mockErr}

	modulePlayers, err := mock.GetModulePlayers(context.Background())

	if !errors.Is(err, mockErr) {
		t.Errorf("expected error %v, got %v", mockErr, err)
	}
	if modulePlayers != nil {
		t.Errorf("expected nil result, got %+v", modulePlayers)
	}
	if mock.NumCalledGetModulePlayers != 1 {
		t.Errorf("expected NumCalledGetModulePlayers 1, got %d", mock.NumCalledGetModulePlayers)
	}
}

func TestMockClient_GetModulePlayers_Reset(t *testing.T) {
	mock := &platform.MockClient{
		NumCalledGetModulePlayers: 3,
		GetModulePlayersError:     errors.New("boom"),
		GetModulePlayersReturn:    []platform.ModulePlayer{{ID: 1}},
	}

	mock.Reset()

	if mock.NumCalledGetModulePlayers != 0 {
		t.Errorf("expected NumCalledGetModulePlayers 0, got %d", mock.NumCalledGetModulePlayers)
	}
	if mock.GetModulePlayersError != nil {
		t.Errorf("expected GetModulePlayersError nil, got %v", mock.GetModulePlayersError)
	}
	if mock.GetModulePlayersReturn != nil {
		t.Errorf("expected GetModulePlayersReturn nil, got %+v", mock.GetModulePlayersReturn)
	}
}
