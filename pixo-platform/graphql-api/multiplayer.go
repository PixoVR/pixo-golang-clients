package graphql_api

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
)

func (g *GraphQLAPIClient) GetMultiplayerServerConfigs(ctx context.Context, params *MultiplayerServerConfigParams) ([]*MultiplayerServerConfigQueryParams, error) {

	variables := map[string]interface{}{
		"params": params,
	}

	var query MultiplayerServerConfigQuery
	if err := g.Client.Query(ctx, &query, variables); err != nil {
		return nil, err
	}

	return query.MultiplayerServerConfigs, nil
}

func (g *GraphQLAPIClient) GetMultiplayerServerVersions(ctx context.Context, params *MultiplayerServerVersionQueryParams) ([]*MultiplayerServerVersion, error) {

	configs, err := g.GetMultiplayerServerConfigs(ctx, &MultiplayerServerConfigParams{
		ModuleID:      params.ModuleID,
		ServerVersion: params.SemanticVersion,
	})
	if err != nil {
		return nil, err
	}

	if len(configs) == 0 {
		return nil, errors.New("no multiplayer server configurations found")
	}

	res := make([]*MultiplayerServerVersion, len(configs[0].ServerVersions))

	for i := range configs[0].ServerVersions {
		res[i] = &MultiplayerServerVersion{
			ModuleID:        configs[0].ModuleID,
			ImageRegistry:   configs[0].ServerVersions[i].ImageRegistry,
			SemanticVersion: configs[0].ServerVersions[i].SemanticVersion,
		}
	}

	return res, nil
}

func (g *GraphQLAPIClient) GetMultiplayerServerVersion(ctx context.Context, versionID int) (*MultiplayerServerVersion, error) {
	query := `query multiplayerServerVersion($id: ID!) { multiplayerServerVersion(id: $id) { id moduleId imageRegistry engine status semanticVersion module { name } } }`

	variables := map[string]interface{}{
		"id": versionID,
	}

	res, err := g.Client.ExecRaw(ctx, query, variables)
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

func (g *GraphQLAPIClient) CreateMultiplayerServerVersion(ctx context.Context, input MultiplayerServerVersion) (*MultiplayerServerVersion, error) {
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
		res, err := g.Client.ExecRaw(ctx, query, variables)
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

	g.AbstractServiceClient.SetHeader("Content-Type", writer.FormDataContentType())

	res, err := g.Post("query", payload.Bytes())
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
