package platform

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"strings"
	"testing"

	abstract "github.com/PixoVR/pixo-golang-clients/pixo-platform/abstract"
	"github.com/PixoVR/pixo-golang-clients/pixo-platform/urlfinder"
)

// newStubbedClient serves responseBody to every query and records the query asked for.
func newStubbedClient(t *testing.T, responseBody string) (Client, *string) {
	t.Helper()

	var askedQuery string

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var request GraphQLRequestPayload
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			t.Errorf("could not decode the request: %v", err)
		}
		askedQuery = request.Query

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(responseBody))
	}))
	t.Cleanup(server.Close)

	serverURL, err := url.Parse(server.URL)
	if err != nil {
		t.Fatalf("could not parse the stub server url: %v", err)
	}

	port, err := strconv.Atoi(serverURL.Port())
	if err != nil {
		t.Fatalf("could not read the stub server port: %v", err)
	}

	client := &clientImpl{
		ServiceClient: abstract.NewClient(abstract.Config{
			ServiceConfig: urlfinder.ServiceConfig{
				Service:     "v2",
				ServiceName: "primary-api",
				Lifecycle:   "local",
				Port:        port,
			},
			Token: "test-token",
		}),
	}

	return client, &askedQuery
}

func TestGetModule_AsksThePlayerForItsEnabledVersions(t *testing.T) {
	client, askedQuery := newStubbedClient(t, `{"data":{"module":{"id":1}}}`)

	if _, err := client.GetModule(context.Background(), 1); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	playerSelection := `modulePlayer {
		versions(status: ["enabled"]) {
			id
			modulePlayerId
			status
			version
			fileName
			package
			platforms { id name shortName }
		}
	}`
	if !strings.Contains(*askedQuery, playerSelection) {
		t.Errorf("expected the query to ask the module player for its enabled versions, got %s", *askedQuery)
	}
}

func TestGetModule_ReadsThePlayerVersionsOfTheModule(t *testing.T) {
	client, _ := newStubbedClient(t, `{"data":{"module":{
		"id": 1,
		"abbreviation": "TST",
		"modulePlayer": {
			"id": 2,
			"name": "Pixo Player",
			"versions": [
				{
					"id": 3,
					"modulePlayerId": 2,
					"status": "enabled",
					"version": "1.2.3",
					"fileName": "player.apk",
					"package": "com.pixovr.player",
					"platforms": [{"id": 1, "name": "Meta Quest 2", "shortName": "oculusquest2"}]
				}
			]
		}
	}}}`)

	module, err := client.GetModule(context.Background(), 1)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if module.ModulePlayer == nil {
		t.Fatal("expected the module to carry its player")
	}
	if len(module.ModulePlayer.Versions) != 1 {
		t.Fatalf("expected one player version, got %+v", module.ModulePlayer.Versions)
	}

	version := module.ModulePlayer.Versions[0]
	if version.ID != 3 || version.ModulePlayerID != 2 {
		t.Errorf("expected version 3 of player 2, got %+v", version)
	}
	if version.Version != "1.2.3" || version.Status != "enabled" {
		t.Errorf("expected the enabled version 1.2.3, got %+v", version)
	}
	if version.FileName != "player.apk" || version.Package != "com.pixovr.player" {
		t.Errorf("expected the file name and package of the version, got %+v", version)
	}
	if len(version.Platforms) != 1 || version.Platforms[0].ShortName != "oculusquest2" {
		t.Errorf("expected the version to support the quest, got %+v", version.Platforms)
	}
}
