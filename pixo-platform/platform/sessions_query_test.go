package platform_test

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/PixoVR/pixo-golang-clients/pixo-platform/platform"
	"github.com/PixoVR/pixo-golang-clients/pixo-platform/urlfinder"
)

type recordingTransport struct {
	queries []string
	reply   string
}

func (r *recordingTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	var body struct {
		Query string `json:"query"`
	}
	raw, _ := io.ReadAll(req.Body)
	_ = json.Unmarshal(raw, &body)
	r.queries = append(r.queries, body.Query)
	return &http.Response{
		StatusCode: http.StatusOK,
		Header:     http.Header{"Content-Type": []string{"application/json"}},
		Body:       io.NopCloser(bytes.NewBufferString(r.reply)),
		Request:    req,
	}, nil
}

func recordQueries(t *testing.T, reply string) *recordingTransport {
	t.Helper()
	transport := &recordingTransport{reply: reply}
	original := http.DefaultClient.Transport
	http.DefaultClient.Transport = transport
	t.Cleanup(func() { http.DefaultClient.Transport = original })
	return transport
}

func localClient() platform.Client {
	return platform.NewClient(urlfinder.ClientConfig{Lifecycle: "local", Token: "test-token"})
}

func TestCreateSessionDoesNotRequestTheSessionsEvents(t *testing.T) {
	transport := recordQueries(t, `{"data":{"createSession":{"id":7,"userId":3,"moduleId":1,"user":{"orgId":2},"module":{"id":1,"abbreviation":"MOD"}}}}`)
	session := &platform.Session{ModuleID: 1}

	if err := localClient().CreateSession(context.Background(), session); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(transport.queries) != 1 {
		t.Fatalf("expected 1 request, got %d", len(transport.queries))
	}
	if strings.Contains(transport.queries[0], "events") {
		t.Errorf("createSession should not request the session's events: %s", transport.queries[0])
	}
	if session.ID != 7 || session.UserID != 3 || session.ModuleID != 1 {
		t.Errorf("unexpected session: %+v", session)
	}
	if session.User.OrgID != 2 || session.Module.ID != 1 {
		t.Errorf("expected user and module on the session: %+v", session)
	}
}

func TestUpdateSessionDoesNotRequestTheSessionsEvents(t *testing.T) {
	transport := recordQueries(t, `{"data":{"updateSession":{"id":7,"userId":3,"moduleId":1,"duration":"1m30s","user":{"orgId":2}}}}`)

	updated, err := localClient().UpdateSession(context.Background(), platform.Session{ID: 7, Completed: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(transport.queries) != 1 {
		t.Fatalf("expected 1 request, got %d", len(transport.queries))
	}
	if strings.Contains(transport.queries[0], "events") {
		t.Errorf("updateSession should not request the session's events: %s", transport.queries[0])
	}
	if updated.ID != 7 || updated.Duration != "1m30s" {
		t.Errorf("unexpected session: %+v", updated)
	}
}
