package platform

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/rs/zerolog/log"
	"io"
	"mime"
	"mime/multipart"
	"net/textproto"
	"os"
	"path/filepath"
	"time"
)

type MultiplayerServerConfig struct {
	ID              int    `json:"id"`
	Capacity        int    `json:"capacity,omitempty"`
	StandbyReplicas string `json:"standbyReplicas,omitempty"`
	Disabled        bool   `json:"disabled,omitempty"`

	ModuleID int     `json:"moduleId,omitempty"`
	Module   *Module `json:"module,omitempty"`

	ServerVersions []MultiplayerServerVersion `json:"serverVersions,omitempty"`

	CreatedBy string `json:"createdBy,omitempty"`
	UpdatedBy string `json:"updatedBy,omitempty"`

	CreatedAt *time.Time `json:"createdAt,omitempty"`
	UpdatedAt *time.Time `json:"updatedAt,omitempty"`
}

type MultiplayerServerTrigger struct {
	ID         int    `json:"id,omitempty"`
	Revision   string `json:"revision,omitempty"`
	Dockerfile string `json:"dockerfile,omitempty"`
	Context    string `json:"context,omitempty"`
	Config     string `json:"config,omitempty"`

	Module   *Module `json:"module,omitempty"`
	ModuleID int     `json:"moduleId,omitempty"`

	CreatedBy string `json:"createdBy,omitempty"`
	UpdatedBy string `json:"updatedBy,omitempty"`

	CreatedAt *time.Time `json:"createdAt,omitempty"`
	UpdatedAt *time.Time `json:"updatedAt,omitempty"`
}

type MultiplayerServerVersion struct {
	ID              int    `json:"id,omitempty" graphql:"id"`
	Engine          string `json:"engine,omitempty" graphql:"engine"`
	Status          string `json:"status,omitempty" graphql:"status"`
	ImageRegistry   string `json:"imageRegistry" graphql:"imageRegistry"`
	SemanticVersion string `json:"semanticVersion,omitempty" graphql:"semanticVersion"`
	FileLink        string `json:"fileLink,omitempty" graphql:"fileLink"`
	LocalFilePath   string `json:"-" graphql:"-"`

	ModuleID int     `json:"moduleId,omitempty" graphql:"moduleId"`
	Module   *Module `json:"module,omitempty" graphql:"module"`

	CreatedBy string `json:"createdBy" graphql:"createdBy"`
	UpdatedBy string `json:"updatedBy" graphql:"updatedBy"`

	CreatedAt *time.Time `json:"createdAt,omitempty" graphql:"createdAt"`
	UpdatedAt *time.Time `json:"updatedAt,omitempty" graphql:"updatedAt"`
}

func (p *PlatformClient) GetMultiplayerServerConfigs(ctx context.Context, params *MultiplayerServerConfigParams) ([]MultiplayerServerConfigQueryParams, error) {

	variables := map[string]interface{}{
		"params": params,
	}

	var query MultiplayerServerConfigQuery
	if err := p.Client.Query(ctx, &query, variables); err != nil {
		return nil, err
	}

	return query.MultiplayerServerConfigs, nil
}

func (p *PlatformClient) GetMultiplayerServerVersions(ctx context.Context, params *MultiplayerServerVersionParams) ([]MultiplayerServerVersion, error) {

	variables := map[string]interface{}{
		"params": params,
	}

	var query MultiplayerServerVersionQuery
	if err := p.Client.Query(ctx, &query, variables); err != nil {
		return nil, err
	}

	return query.MultiplayerServerVersions, nil
}

func (p *PlatformClient) GetMultiplayerServerVersionsWithConfig(ctx context.Context, params *MultiplayerServerVersionParams) ([]MultiplayerServerVersion, error) {

	configs, err := p.GetMultiplayerServerConfigs(ctx, &MultiplayerServerConfigParams{
		ModuleID:      params.ModuleID,
		ServerVersion: params.SemanticVersion,
	})
	if err != nil {
		return nil, err
	}

	if len(configs) == 0 {
		return nil, errors.New("no multiplayer server configurations found")
	}

	res := make([]MultiplayerServerVersion, len(configs[0].ServerVersions))

	for i := range configs[0].ServerVersions {
		res[i] = MultiplayerServerVersion{
			ID:              configs[0].ServerVersions[i].ID,
			ModuleID:        configs[0].ModuleID,
			ImageRegistry:   configs[0].ServerVersions[i].ImageRegistry,
			SemanticVersion: configs[0].ServerVersions[i].SemanticVersion,
		}
	}

	return res, nil
}

func (p *PlatformClient) GetMultiplayerServerVersion(ctx context.Context, versionID int) (*MultiplayerServerVersion, error) {
	query := `query multiplayerServerVersion($id: ID!) { multiplayerServerVersion(id: $id) { id moduleId imageRegistry engine status semanticVersion module { name } } }`

	variables := map[string]interface{}{
		"id": versionID,
	}

	res, err := p.Client.ExecRaw(ctx, query, variables)
	if err != nil {
		return nil, err
	}

	var response struct {
		MultiplayerServerVersion *MultiplayerServerVersion `json:"multiplayerServerVersion"`
	}
	if err = json.Unmarshal(res, &response); err != nil {
		return nil, err
	}

	return response.MultiplayerServerVersion, nil
}

func (p *PlatformClient) CreateMultiplayerServerVersion(ctx context.Context, input MultiplayerServerVersion) (*MultiplayerServerVersion, error) {
	query := `mutation createMultiplayerServerVersion($input: MultiplayerServerVersionInput!) { createMultiplayerServerVersion(input: $input) { id imageRegistry fileLink semanticVersion engine module { name } } }`

	if input.ImageRegistry == "" && input.LocalFilePath == "" {
		return nil, errors.New("image or file path must be provided")
	}

	if input.Status == "" {
		input.Status = "enabled"
	}

	variables := map[string]interface{}{
		"input": map[string]interface{}{
			"moduleId":        input.ModuleID,
			"imageRegistry":   input.ImageRegistry,
			"semanticVersion": input.SemanticVersion,
			"engine":          input.Engine,
			"status":          input.Status,
		},
	}

	if input.LocalFilePath == "" {
		res, err := p.Client.ExecRaw(ctx, query, variables)
		if err != nil {
			return nil, err
		}

		var response struct {
			ServerVersion *MultiplayerServerVersion `json:"createMultiplayerServerVersion"`
		}
		if err = json.Unmarshal(res, &response); err != nil {
			return nil, err
		}

		return response.ServerVersion, nil
	}

	graphqlRequest := struct {
		OperationName string         `json:"operationName"`
		Query         string         `json:"query"`
		Variables     map[string]any `json:"variables,omitempty"`
	}{
		OperationName: "createMultiplayerServerVersion",
		Query:         query,
		Variables:     variables,
	}

	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(graphqlRequest); err != nil {
		return nil, err
	}

	payload := &bytes.Buffer{}
	writer := multipart.NewWriter(payload)
	_ = writer.WriteField("operations", buf.String())
	file, err := os.Open(input.LocalFilePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	mapData := map[string][]string{}
	mapData["0"] = []string{fmt.Sprintf(`variables.%s`, "input.filePath")}
	jsonData, _ := json.Marshal(mapData)

	_ = writer.WriteField("map", string(jsonData))

	part, err := createFormFile(writer, "0", filepath.Base(input.LocalFilePath))
	if err != nil {
		return nil, err
	}

	if _, err = io.Copy(part, file); err != nil {
		return nil, err
	}

	if err = writer.Close(); err != nil {
		return nil, err
	}

	p.AbstractServiceClient.SetHeader("Content-Type", writer.FormDataContentType())

	res, err := p.Post("query", payload.Bytes())
	if err != nil {
		log.Error().Err(err).Msg("error creating multiplayer server version")
		return nil, err
	}

	if res.IsError() {
		return nil, fmt.Errorf("error creating multiplayer server version: %s", res.String())
	}

	var gqlRes struct {
		Data struct {
			CreateMultiplayerServerVersion *MultiplayerServerVersion `json:"createMultiplayerServerVersion"`
		} `json:"data"`
	}

	if err = json.Unmarshal(res.Body(), &gqlRes); err != nil {
		return nil, err
	}

	return gqlRes.Data.CreateMultiplayerServerVersion, nil
}

func createFormFile(w *multipart.Writer, fieldName, filename string) (io.Writer, error) {
	fileContentType := mime.TypeByExtension(filepath.Ext(filename))
	if fileContentType == "" {
		fileContentType = "application/octet-stream"
	}
	h := make(textproto.MIMEHeader)
	h.Set("Content-Disposition", fmt.Sprintf(`form-data; name="%s"; filename="%s"`, fieldName, filename))
	h.Set("Content-Type", fileContentType)
	return w.CreatePart(h)
}
