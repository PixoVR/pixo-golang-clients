package platform_test

import (
	"context"
	"errors"
	"testing"

	"github.com/PixoVR/pixo-golang-clients/pixo-platform/platform"
)

func TestMockClient_GetModulePlayerVersions_Default(t *testing.T) {
	mock := &platform.MockClient{}

	versions, err := mock.GetModulePlayerVersions(context.Background(), nil)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(versions) != 1 {
		t.Fatalf("expected 1 module player version, got %d", len(versions))
	}
	got := versions[0]
	if got.ID != 1 {
		t.Errorf("expected ID 1, got %d", got.ID)
	}
	if got.ModulePlayerID != 1 {
		t.Errorf("expected ModulePlayerID 1, got %d", got.ModulePlayerID)
	}
	if got.Version != "1.0.0" {
		t.Errorf("expected Version 1.0.0, got %q", got.Version)
	}
	if got.Status != "enabled" {
		t.Errorf("expected Status enabled, got %q", got.Status)
	}
	if mock.NumCalledGetModulePlayerVersions != 1 {
		t.Errorf("expected NumCalledGetModulePlayerVersions 1, got %d", mock.NumCalledGetModulePlayerVersions)
	}
}

func TestMockClient_GetModulePlayerVersions_CapturesParams(t *testing.T) {
	mock := &platform.MockClient{}
	id := 42
	params := &platform.ModulePlayerVersionParams{
		Status:         []string{"enabled"},
		ModulePlayerID: &id,
	}

	_, err := mock.GetModulePlayerVersions(context.Background(), params)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if mock.GetModulePlayerVersionsParameters != params {
		t.Errorf("expected params to be captured, got %+v", mock.GetModulePlayerVersionsParameters)
	}
}

func TestMockClient_GetModulePlayerVersions_ReturnsConfigured(t *testing.T) {
	expected := []platform.ModulePlayerVersion{
		{ID: 7, ModulePlayerID: 3, Version: "2.0.0", Status: "disabled"},
	}
	mock := &platform.MockClient{GetModulePlayerVersionsReturn: expected}

	versions, err := mock.GetModulePlayerVersions(context.Background(), nil)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(versions) != 1 || versions[0].ID != 7 || versions[0].Version != "2.0.0" {
		t.Errorf("unexpected result: %+v", versions)
	}
}

func TestMockClient_GetModulePlayerVersions_ReturnsError(t *testing.T) {
	mockErr := errors.New("boom")
	mock := &platform.MockClient{GetModulePlayerVersionsError: mockErr}

	versions, err := mock.GetModulePlayerVersions(context.Background(), nil)

	if !errors.Is(err, mockErr) {
		t.Errorf("expected error %v, got %v", mockErr, err)
	}
	if versions != nil {
		t.Errorf("expected nil result, got %+v", versions)
	}
}

func TestMockClient_GetModulePlayerVersions_Reset(t *testing.T) {
	mock := &platform.MockClient{
		NumCalledGetModulePlayerVersions:  3,
		GetModulePlayerVersionsError:      errors.New("boom"),
		GetModulePlayerVersionsParameters: &platform.ModulePlayerVersionParams{},
		GetModulePlayerVersionsReturn:     []platform.ModulePlayerVersion{{ID: 1}},
	}

	mock.Reset()

	if mock.NumCalledGetModulePlayerVersions != 0 {
		t.Errorf("expected NumCalledGetModulePlayerVersions 0, got %d", mock.NumCalledGetModulePlayerVersions)
	}
	if mock.GetModulePlayerVersionsError != nil {
		t.Errorf("expected GetModulePlayerVersionsError nil, got %v", mock.GetModulePlayerVersionsError)
	}
	if mock.GetModulePlayerVersionsParameters != nil {
		t.Errorf("expected GetModulePlayerVersionsParameters nil, got %+v", mock.GetModulePlayerVersionsParameters)
	}
	if mock.GetModulePlayerVersionsReturn != nil {
		t.Errorf("expected GetModulePlayerVersionsReturn nil, got %+v", mock.GetModulePlayerVersionsReturn)
	}
}
