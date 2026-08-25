package platform

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"strings"
	"testing"

	abstract "github.com/PixoVR/pixo-golang-clients/pixo-platform/abstract"
	"github.com/PixoVR/pixo-golang-clients/pixo-platform/urlfinder"
)

// newStubbedClient serves responseBody to every query and records the request made.
func newStubbedClient(t *testing.T, responseBody string) (Client, *GraphQLRequestPayload) {
	t.Helper()

	var asked GraphQLRequestPayload

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var request GraphQLRequestPayload
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			t.Errorf("could not decode the request: %v", err)
		}
		asked = request

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

	return client, &asked
}

func TestGetModule_AsksThePlayerForItsEnabledVersions(t *testing.T) {
	client, asked := newStubbedClient(t, `{"data":{"module":{"id":1}}}`)

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
	if !strings.Contains(asked.Query, playerSelection) {
		t.Errorf("expected the query to ask the module player for its enabled versions, got %s", asked.Query)
	}
}

func TestGetModulesWithAssociations_AsksForTheNamedModulesWithTheirVersions(t *testing.T) {
	client, asked := newStubbedClient(t, `{"data":{"modules":[]}}`)

	_, err := client.GetModulesWithAssociations(context.Background(), ModuleParams{
		IDs:          []int{11, 22},
		LifecycleIds: []int{3},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(asked.Query, "modules(params: $params)") {
		t.Errorf("expected one modules query, got %s", asked.Query)
	}
	if !strings.Contains(asked.Query, "versions {") {
		t.Errorf("expected the query to ask each module for its versions, got %s", asked.Query)
	}
	if !strings.Contains(asked.Query, `versions(status: ["enabled"])`) {
		t.Errorf("expected the query to ask each module player for its enabled versions, got %s", asked.Query)
	}

	params, ok := asked.Variables["params"].(map[string]interface{})
	if !ok {
		t.Fatalf("expected the query to send filter params, got %+v", asked.Variables)
	}
	if fmt.Sprint(params["ids"]) != "[11 22]" {
		t.Errorf("expected the query to name modules 11 and 22, got %+v", params["ids"])
	}
	if fmt.Sprint(params["lifecycleIds"]) != "[3]" {
		t.Errorf("expected the query to ask for the released lifecycle, got %+v", params["lifecycleIds"])
	}
}

func TestGetModulesWithAssociations_ReadsTheVersionsAndPlayerVersionsOfEachModule(t *testing.T) {
	client, _ := newStubbedClient(t, `{"data":{"modules":[
		{
			"id": 11,
			"abbreviation": "ONE",
			"versions": [
				{
					"id": 101,
					"moduleId": 11,
					"semanticVersion": "1.0.0",
					"filePath": "ModuleVersions/101/zips/module.zip",
					"platforms": [{"id": 1, "shortName": "webgl"}]
				}
			],
			"modulePlayer": {
				"id": 2,
				"versions": [{"id": 3, "modulePlayerId": 2, "status": "enabled", "version": "1.2.3"}]
			}
		},
		{"id": 22, "abbreviation": "TWO"}
	]}}`)

	modules, err := client.GetModulesWithAssociations(context.Background(), ModuleParams{IDs: []int{11, 22}})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(modules) != 2 {
		t.Fatalf("expected both modules, got %+v", modules)
	}

	if modules[0].ID != 11 || modules[1].ID != 22 {
		t.Errorf("expected modules 11 and 22, got %d and %d", modules[0].ID, modules[1].ID)
	}
	if len(modules[0].Versions) != 1 || modules[0].Versions[0].ID != 101 {
		t.Fatalf("expected module 11 to carry version 101, got %+v", modules[0].Versions)
	}
	if modules[0].Versions[0].FilePath != "ModuleVersions/101/zips/module.zip" {
		t.Errorf("expected the storage path of version 101, got %+v", modules[0].Versions[0])
	}
	if len(modules[0].Versions[0].Platforms) != 1 || modules[0].Versions[0].Platforms[0].ShortName != "webgl" {
		t.Errorf("expected version 101 to support webgl, got %+v", modules[0].Versions[0].Platforms)
	}
	if modules[0].ModulePlayer == nil || len(modules[0].ModulePlayer.Versions) != 1 {
		t.Fatalf("expected module 11 to carry the enabled versions of its player, got %+v", modules[0].ModulePlayer)
	}
	if modules[0].ModulePlayer.Versions[0].Version != "1.2.3" {
		t.Errorf("expected player version 1.2.3, got %+v", modules[0].ModulePlayer.Versions[0])
	}
	if len(modules[1].Versions) != 0 {
		t.Errorf("expected module 22 to carry no version, got %+v", modules[1].Versions)
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
