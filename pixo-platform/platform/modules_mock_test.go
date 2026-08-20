package platform_test

import (
	"context"
	"errors"
	"testing"

	"github.com/PixoVR/pixo-golang-clients/pixo-platform/platform"
)

func TestMockClient_GetModules_CapturesParams(t *testing.T) {
	mock := &platform.MockClient{}
	isPublic := true
	params := platform.ModuleParams{IsPublic: &isPublic, Statuses: []string{"enabled"}}

	_, err := mock.GetModules(context.Background(), params)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if mock.GetModulesParameters == nil {
		t.Fatal("expected params to be captured")
	}
	if mock.GetModulesParameters.IsPublic != &isPublic {
		t.Errorf("expected IsPublic to be captured, got %+v", mock.GetModulesParameters.IsPublic)
	}
}

func TestMockClient_GetModules_ReturnsConfigured(t *testing.T) {
	expected := []platform.Module{{ID: 42, Abbreviation: "CUSTOM"}}
	mock := &platform.MockClient{GetModulesReturn: expected}

	modules, err := mock.GetModules(context.Background())

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(modules) != 1 || modules[0].ID != 42 {
		t.Errorf("unexpected result: %+v", modules)
	}
}

func TestMockClient_GetModule_Default(t *testing.T) {
	mock := &platform.MockClient{}

	module, err := mock.GetModule(context.Background(), 7)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if module.ID != 7 {
		t.Errorf("expected ID 7, got %d", module.ID)
	}
	if len(module.Versions) != 1 || module.Versions[0].ModuleID != 7 {
		t.Errorf("expected one version belonging to the module, got %+v", module.Versions)
	}
	if mock.GetModuleIDParam != 7 {
		t.Errorf("expected GetModuleIDParam 7, got %d", mock.GetModuleIDParam)
	}
	if mock.NumCalledGetModule != 1 {
		t.Errorf("expected NumCalledGetModule 1, got %d", mock.NumCalledGetModule)
	}
}

func TestMockClient_GetModule_RequiresID(t *testing.T) {
	mock := &platform.MockClient{}

	module, err := mock.GetModule(context.Background(), 0)

	if err == nil {
		t.Error("expected an error when the module id is missing")
	}
	if module != nil {
		t.Errorf("expected nil result, got %+v", module)
	}
}

func TestMockClient_GetModule_ReturnsError(t *testing.T) {
	mockErr := errors.New("boom")
	mock := &platform.MockClient{GetModuleError: mockErr}

	module, err := mock.GetModule(context.Background(), 1)

	if !errors.Is(err, mockErr) {
		t.Errorf("expected error %v, got %v", mockErr, err)
	}
	if module != nil {
		t.Errorf("expected nil result, got %+v", module)
	}
}

func TestMockClient_GetOrgModules_Default(t *testing.T) {
	mock := &platform.MockClient{}

	orgModules, err := mock.GetOrgModules(context.Background(), platform.OrgModuleParams{OrgID: 3})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(orgModules) != 1 || orgModules[0].OrgID != 3 {
		t.Fatalf("expected one org module for org 3, got %+v", orgModules)
	}
	if orgModules[0].Module == nil {
		t.Error("expected the org module to include its module")
	}
	if mock.GetOrgModulesParameters == nil || mock.GetOrgModulesParameters.OrgID != 3 {
		t.Errorf("expected params to be captured, got %+v", mock.GetOrgModulesParameters)
	}
}

func TestMockClient_GetOrgModules_RequiresOrgID(t *testing.T) {
	mock := &platform.MockClient{}

	orgModules, err := mock.GetOrgModules(context.Background(), platform.OrgModuleParams{})

	if err == nil {
		t.Error("expected an error when the org id is missing")
	}
	if orgModules != nil {
		t.Errorf("expected nil result, got %+v", orgModules)
	}
}

func TestMockClient_GetOrgModules_ReturnsError(t *testing.T) {
	mockErr := errors.New("boom")
	mock := &platform.MockClient{GetOrgModulesError: mockErr}

	orgModules, err := mock.GetOrgModules(context.Background(), platform.OrgModuleParams{OrgID: 1})

	if !errors.Is(err, mockErr) {
		t.Errorf("expected error %v, got %v", mockErr, err)
	}
	if orgModules != nil {
		t.Errorf("expected nil result, got %+v", orgModules)
	}
}

func TestMockClient_ModuleReads_Reset(t *testing.T) {
	mock := &platform.MockClient{
		NumCalledGetModules:     1,
		GetModulesReturn:        []platform.Module{{ID: 1}},
		NumCalledGetModule:      2,
		GetModuleIDParam:        9,
		GetModuleReturn:         &platform.Module{ID: 9},
		NumCalledGetOrgModules:  3,
		GetOrgModulesParameters: &platform.OrgModuleParams{OrgID: 4},
		GetOrgModulesReturn:     []platform.OrgModule{{ID: 1}},
	}

	mock.Reset()

	if mock.NumCalledGetModules != 0 || mock.GetModulesReturn != nil {
		t.Error("expected the GetModules state to be reset")
	}
	if mock.NumCalledGetModule != 0 || mock.GetModuleIDParam != 0 || mock.GetModuleReturn != nil {
		t.Error("expected the GetModule state to be reset")
	}
	if mock.NumCalledGetOrgModules != 0 || mock.GetOrgModulesParameters != nil || mock.GetOrgModulesReturn != nil {
		t.Error("expected the GetOrgModules state to be reset")
	}
}

func TestMockClient_GetUsersWithModuleAccess_CapturesParams(t *testing.T) {
	mock := &platform.MockClient{}

	users, err := mock.GetUsersWithModuleAccess(context.Background(), 43, 20)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if mock.NumCalledGetUsersWithModuleAccess != 1 {
		t.Errorf("expected one call, got %d", mock.NumCalledGetUsersWithModuleAccess)
	}
	if mock.GetUsersWithModuleAccessModuleID != 43 || mock.GetUsersWithModuleAccessOrgID != 20 {
		t.Errorf("expected the module and org to be captured, got %d and %d", mock.GetUsersWithModuleAccessModuleID, mock.GetUsersWithModuleAccessOrgID)
	}
	if len(users) != 1 || users[0].OrgID != 20 {
		t.Errorf("unexpected result: %+v", users)
	}
}

func TestMockClient_GetUsersWithModuleAccess_ReturnsError(t *testing.T) {
	mock := &platform.MockClient{GetUsersWithModuleAccessError: errors.New("boom")}

	if _, err := mock.GetUsersWithModuleAccess(context.Background(), 43, 20); err == nil {
		t.Fatal("expected an error")
	}
}
