package platform

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"strconv"
	"testing"

	"github.com/PixoVR/pixo-golang-clients/pixo-platform/abstract"
	"github.com/PixoVR/pixo-golang-clients/pixo-platform/urlfinder"
)

func authClient(t *testing.T, status int, body string) (Client, *httptest.Server) {
	t.Helper()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet || r.URL.Path != "/v2/auth/check" {
			t.Errorf("unexpected request: %s %s", r.Method, r.URL.Path)
		}
		w.WriteHeader(status)
		_, _ = w.Write([]byte(body))
	}))
	t.Cleanup(server.Close)
	serverURL, err := url.Parse(server.URL)
	if err != nil {
		t.Fatal(err)
	}
	port, err := strconv.Atoi(serverURL.Port())
	if err != nil {
		t.Fatal(err)
	}
	return &clientImpl{ServiceClient: abstract.NewClient(abstract.Config{
		ServiceConfig: urlfinder.ServiceConfig{Service: "v2", ServiceName: "primary-api", Lifecycle: "local", Port: port},
		Token:         "test-token",
	})}, server
}

func TestFeaturePermissions(t *testing.T) {
	for _, tc := range []struct {
		name string
		body string
		want []string
	}{
		{"granted", `{"user":{"featurePermissions":["web-module-service:manage"]}}`, []string{"web-module-service:manage"}},
		{"empty", `{"user":{"featurePermissions":[]}}`, []string{}},
		{"omitted", `{"user":{}}`, []string{}},
		{"null", `{"user":{"featurePermissions":null}}`, []string{}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			client, _ := authClient(t, http.StatusOK, tc.body)
			got, err := client.GetFeaturePermissions(context.Background())
			if err != nil || !reflect.DeepEqual(got, tc.want) {
				t.Fatalf("got %v, %v; want %v", got, err, tc.want)
			}
			for _, permission := range []string{"web-module-service:manage", "other"} {
				got, err := client.HasFeaturePermission(context.Background(), permission)
				want := len(tc.want) > 0 && permission == tc.want[0]
				if err != nil || got != want {
					t.Fatalf("permission %q: got %v, %v; want %v", permission, got, err, want)
				}
			}
		})
	}
}

func TestFeaturePermissionErrors(t *testing.T) {
	for _, tc := range []struct {
		name         string
		status       int
		body         string
		unauthorized bool
	}{
		{"bad request", 400, `{"error":"no token provided"}`, true},
		{"unauthorized", 401, `{"error":"unauthorized"}`, true},
		{"forbidden", 403, `{"error":"forbidden"}`, true},
		{"empty unauthorized", 401, "", true},
		{"server error", 500, `{"error":"unavailable"}`, false},
		{"non JSON server error", 502, "bad gateway", false},
		{"invalid success payload", 200, "invalid", false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			client, _ := authClient(t, tc.status, tc.body)
			_, err := client.CheckAuth(context.Background())
			if err == nil || errors.Is(err, ErrUnauthorized) != tc.unauthorized {
				t.Fatalf("CheckAuth: got %v; unauthorized=%v", err, tc.unauthorized)
			}
			permissions, err := client.GetFeaturePermissions(context.Background())
			if err == nil || permissions != nil || errors.Is(err, ErrUnauthorized) != tc.unauthorized {
				t.Fatalf("GetFeaturePermissions: got %v, %v", permissions, err)
			}
			has, err := client.HasFeaturePermission(context.Background(), "web-module-service:manage")
			if err == nil || has || errors.Is(err, ErrUnauthorized) != tc.unauthorized {
				t.Fatalf("HasFeaturePermission: got %v, %v", has, err)
			}
		})
	}
}

func TestFeaturePermissionsConnectionFailure(t *testing.T) {
	client, server := authClient(t, http.StatusOK, `{}`)
	server.Close()
	_, err := client.GetFeaturePermissions(context.Background())
	if err == nil || errors.Is(err, ErrUnauthorized) {
		t.Fatalf("expected a non-authorization connection error, got %v", err)
	}
}

func TestMockFeaturePermissions(t *testing.T) {
	ctx := context.Background()
	client := &MockClient{GetFeaturePermissionsReturn: []string{"web-module-service:manage"}}
	has, err := client.HasFeaturePermission(ctx, "web-module-service:manage")
	if err != nil || !has {
		t.Fatalf("expected permission, got %v, %v", has, err)
	}
	has, err = client.HasFeaturePermission(ctx, "other")
	if err != nil || has || client.NumCalledGetFeaturePermissions != 2 {
		t.Fatalf("unexpected missing permission result: %v, %v", has, err)
	}
	client.GetFeaturePermissionsError = ErrUnauthorized
	has, err = client.HasFeaturePermission(ctx, "web-module-service:manage")
	if has || !errors.Is(err, ErrUnauthorized) {
		t.Fatalf("expected authorization failure, got %v, %v", has, err)
	}
	client.Reset()
	if client.NumCalledGetFeaturePermissions != 0 || client.GetFeaturePermissionsError != nil || client.GetFeaturePermissionsReturn != nil {
		t.Fatal("Reset did not clear feature permissions")
	}
	permissions, err := client.GetFeaturePermissions(ctx)
	if err != nil || permissions == nil || len(permissions) != 0 {
		t.Fatalf("expected non-nil empty permissions, got %v, %v", permissions, err)
	}
}
