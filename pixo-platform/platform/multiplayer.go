package platform

import (
	"context"
	"errors"
	"time"
)

type MultiplayerServerConfigQuery struct {
	MultiplayerServerConfigs []MultiplayerServerConfigQueryParams `graphql:"multiplayerServerConfigs(params: $params)"`
}

type MultiplayerServerVersionQuery struct {
	MultiplayerServerVersions []MultiplayerServerVersion `graphql:"multiplayerServerVersions(params: $params)"`
}

type MultiplayerServerConfigParams struct {
	ModuleID       int                        `json:"moduleId,omitempty"`
	OrgID          int                        `json:"orgId,omitempty"`
	ServerVersion  string                     `json:"serverVersion,omitempty"`
	Capacity       int                        `json:"capacity,omitempty"`
	ServerVersions []MultiplayerServerVersion `json:"serverVersions,omitempty"`
}

type MultiplayerServerConfigQueryParams struct {
	ID       int  `json:"id" graphql:"id"`
	ModuleID int  `json:"moduleId" graphql:"moduleId"`
	Capacity int  `json:"capacity" graphql:"capacity"`
	Disabled bool `json:"disabled" graphql:"disabled"`
	Module   struct {
		ID   int    `json:"id" graphql:"id"`
		Name string `json:"name" graphql:"name"`
	}
	ServerVersions []MultiplayerServerVersion `json:"serverVersions" graphql:"serverVersions"`
}

type MultiplayerServerVersionParams struct {
	ModuleID        int    `json:"moduleId" graphql:"moduleId"`
	SemanticVersion string `json:"semanticVersion" graphql:"semanticVersion"`
}

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
	ImageRegistry   string `json:"imageRegistry,omitempty" graphql:"imageRegistry"`
	SemanticVersion string `json:"semanticVersion,omitempty" graphql:"semanticVersion"`
	FileLink        string `json:"fileLink,omitempty" graphql:"fileLink"`
	FilePath        string `json:"filePath,omitempty" graphql:"filePath"`
	LocalFilePath   string `json:"-" graphql:"-"`

	ModuleID int     `json:"moduleId,omitempty" graphql:"moduleId"`
	Module   *Module `json:"module,omitempty" graphql:"module"`

	CreatedBy string `json:"createdBy,omitempty" graphql:"createdBy"`
	UpdatedBy string `json:"updatedBy,omitempty" graphql:"updatedBy"`

	CreatedAt *time.Time `json:"createdAt,omitempty" graphql:"createdAt"`
	UpdatedAt *time.Time `json:"updatedAt,omitempty" graphql:"updatedAt"`
}

func (p *clientImpl) GetMultiplayerServerConfigs(ctx context.Context, params *MultiplayerServerConfigParams) ([]MultiplayerServerConfigQueryParams, error) {
	query := `query multiplayerServerConfigs($params: MultiplayerServerConfigParams) { multiplayerServerConfigs(params: $params) { id moduleId capacity disabled module { id name } serverVersions { id semanticVersion imageRegistry } } }`

	variables := map[string]interface{}{
		"params": params,
	}

	var res MultiplayerServerConfigQuery
	if err := p.Exec(ctx, query, &res, variables); err != nil {
		return nil, err
	}

	return res.MultiplayerServerConfigs, nil
}

func (p *clientImpl) GetMultiplayerServerVersions(ctx context.Context, params *MultiplayerServerVersionParams) ([]MultiplayerServerVersion, error) {
	query := `query multiplayerServerVersions($params: MultiplayerServerVersionParams) { multiplayerServerVersions(params: $params) { id moduleId imageRegistry engine status semanticVersion filePath module { name } } }`

	variables := map[string]interface{}{
		"params": params,
	}

	var res MultiplayerServerVersionQuery
	if err := p.Exec(ctx, query, &res, variables); err != nil {
		return nil, err
	}

	return res.MultiplayerServerVersions, nil
}

func (p *clientImpl) GetMultiplayerServerVersionsWithConfig(ctx context.Context, params *MultiplayerServerVersionParams) ([]MultiplayerServerVersion, error) {

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

func (p *clientImpl) UpdateMultiplayerServerVersion(ctx context.Context, input MultiplayerServerVersion) (*MultiplayerServerVersion, error) {
	query := `mutation updateMultiplayerServerVersion($input: MultiplayerServerVersionInput!) { updateMultiplayerServerVersion(input: $input) { id moduleId imageRegistry engine status semanticVersion module { name } } }`

	variables := map[string]interface{}{
		"input": map[string]interface{}{
			"id":              input.ID,
			"moduleId":        input.ModuleID,
			"semanticVersion": input.SemanticVersion,
			"imageRegistry":   input.ImageRegistry,
		},
	}

	if input.Status != "" {
		variables["input"].(map[string]interface{})["status"] = input.Status
	}

	if input.Engine != "" {
		variables["input"].(map[string]interface{})["engine"] = input.Engine
	}

	var res struct {
		ServerVersion *MultiplayerServerVersion `json:"updateMultiplayerServerVersion"`
	}

	if err := p.Exec(ctx, query, &res, variables); err != nil {
		return nil, err
	}

	return res.ServerVersion, nil
}

func (p *clientImpl) GetMultiplayerServerVersion(ctx context.Context, versionID int) (*MultiplayerServerVersion, error) {
	query := `query multiplayerServerVersion($id: ID!) { multiplayerServerVersion(id: $id) { id moduleId imageRegistry engine status semanticVersion filePath module { name } } }`

	variables := map[string]interface{}{
		"id": versionID,
	}

	var res struct {
		MultiplayerServerVersion *MultiplayerServerVersion `json:"multiplayerServerVersion"`
	}
	if err := p.Exec(ctx, query, &res, variables); err != nil {
		return nil, err
	}

	return res.MultiplayerServerVersion, nil
}

func (p *clientImpl) CreateMultiplayerServerVersion(ctx context.Context, input MultiplayerServerVersion) (*MultiplayerServerVersion, error) {
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

	var res struct {
		ServerVersion *MultiplayerServerVersion `json:"createMultiplayerServerVersion"`
	}

	if err := p.ExecWithFile(ctx, query, &res, variables, input.LocalFilePath, "filePath"); err != nil {
		return nil, err
	}

	return res.ServerVersion, nil
}
