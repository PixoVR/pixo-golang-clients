package platform

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/PixoVR/pixo-golang-clients/pixo-platform/legacy"
	"github.com/rs/zerolog/log"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"time"
)

type GitConfig struct {
	Provider string `json:"provider,omitempty"`
	OrgName  string `json:"orgName,omitempty"`
	RepoName string `json:"repoName,omitempty"`
}

type Module struct {
	ID          int    `json:"id,omitempty"`
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
	ImageLink   string `json:"imageLink,omitempty"`
	ShortDesc   string `json:"shortDesc,omitempty"`

	GitConfigID int       `json:"gitConfigId,omitempty"`
	GitConfig   GitConfig `json:"gitConfig,omitempty"`

	CreatedAt time.Time `json:"createdAt,omitempty"`
	UpdatedAt time.Time `json:"updatedAt,omitempty"`
}

type ModuleVersion struct {
	ID              int           `json:"id,omitempty"`
	ModuleID        int           `json:"moduleId,omitempty"`
	Module          legacy.Module `json:"module,omitempty"`
	Status          string        `json:"status,omitempty"`
	FileLink        string        `json:"fileLink,omitempty"`
	SemanticVersion string        `json:"version,omitempty"`
	Notes           string        `json:"notes,omitempty"`
	Package         string        `json:"package,omitempty"`
	ExternalID      string        `json:"externalId,omitempty"`
	LocalFilePath   string        `json:"-"`
	ControlIds      []int         `json:"controlIds,omitempty"`
	PlatformIds     []int         `json:"platformIds,omitempty"`
}

type ModuleParams struct {
	Name string `json:"name"`
}

type GetModulesResponse struct {
	Modules []Module `json:"modules"`
}

type CreateModuleResponse struct {
	Module Module `json:"createModule"`
}

type CreateModuleVersionResponse struct {
	ModuleVersion ModuleVersion `json:"createModuleVersion"`
}

func (p *PlatformClient) GetModules(ctx context.Context, params ...ModuleParams) ([]Module, error) {
	query := `query modules { modules { id name description imageLink shortDesc gitConfigId gitConfig { provider orgName repoName } createdAt updatedAt } }`

	res, err := p.Client.ExecRaw(ctx, query, nil)
	if err != nil {
		return nil, err
	}

	var response GetModulesResponse
	if err = json.Unmarshal(res, &response); err != nil {
		return nil, err
	}

	return response.Modules, nil
}

func (p *PlatformClient) CreateModuleVersion(ctx context.Context, input ModuleVersion) (*ModuleVersion, error) {
	query := `mutation createModuleVersion($input: ModuleVersionInput!) { createModuleVersion(input: $input) { id moduleId module { abbreviation } version package status fileLink } }`

	if input.LocalFilePath == "" {
		return nil, errors.New("file path must be provided")
	}

	if input.Status == "" {
		input.Status = "disabled"
	}

	variables := map[string]interface{}{
		"input": map[string]interface{}{
			"moduleId":    input.ModuleID,
			"version":     input.SemanticVersion,
			"notes":       input.Notes,
			"status":      input.Status,
			"package":     input.Package,
			"platformIds": input.PlatformIds,
			"controlIds":  input.ControlIds,
		},
	}

	if input.LocalFilePath == "" {
		res, err := p.Client.ExecRaw(ctx, query, variables)
		if err != nil {
			return nil, err
		}

		var response CreateModuleVersionResponse
		if err = json.Unmarshal(res, &response); err != nil {
			return nil, err
		}

		return &response.ModuleVersion, nil
	}

	graphqlRequest := struct {
		OperationName string         `json:"operationName"`
		Query         string         `json:"query"`
		Variables     map[string]any `json:"variables,omitempty"`
	}{
		OperationName: "createModuleVersion",
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
		Data   CreateModuleVersionResponse `json:"data"`
		Errors []struct {
			Message string `json:"message"`
		}
	}

	if err = json.Unmarshal(res.Body(), &gqlRes); err != nil {
		return nil, err
	}

	if len(gqlRes.Errors) > 0 {
		return nil, errors.New(gqlRes.Errors[0].Message)
	}

	return &gqlRes.Data.ModuleVersion, nil
}
