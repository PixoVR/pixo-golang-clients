package platform

import (
	"context"
	"errors"
	"time"
)

type Asset struct {
	ID   int    `json:"id,omitempty"`
	Name string `json:"name,omitempty" yaml:"name,omitempty"`
	Type string `json:"type,omitempty"`

	ModuleID int            `json:"moduleId,omitempty"`
	Module   *Module        `json:"module,omitempty"`
	Versions []AssetVersion `json:"versions,omitempty" yaml:"versions,omitempty"`

	CreatedAt time.Time `json:"createdAt,omitempty"`
	UpdatedAt time.Time `json:"updatedAt,omitempty"`
}

type AssetParams struct {
	ModuleID         int    `json:"moduleId,omitempty"`
	ExternalModuleID string `json:"externalModuleId,omitempty"`
	Name             string `json:"name,omitempty"`
	Type             string `json:"type,omitempty"`
	Status           string `json:"status,omitempty"`
	LanguageCode     string `json:"languageCode,omitempty"`
	Tags             string `json:"tags,omitempty"`
}

type CreateAssetResponse struct {
	Asset Asset `json:"createAsset"`
}

type GetAssetResponse struct {
	Asset Asset `json:"asset"`
}

type GetAssetsResponse struct {
	Assets []Asset `json:"assets"`
}

type AssetVersion struct {
	ID            int    `json:"id,omitempty"`
	Status        string `json:"status,omitempty"`
	LanguageCode  string `json:"languageCode,omitempty"`
	Language      string `json:"language,omitempty"`
	FileLink      string `json:"fileLink,omitempty"`
	LocalFilePath string `json:"-"`

	AssetID int    `json:"assetId,omitempty"`
	Asset   *Asset `json:"asset,omitempty"`

	CreatedAt time.Time `json:"createdAt,omitempty"`
	UpdatedAt time.Time `json:"updatedAt,omitempty"`
}

type CreateAssetVersionResponse struct {
	AssetVersion AssetVersion `json:"createAssetVersion"`
}

type UpdateAssetVersionResponse struct {
	AssetVersion AssetVersion `json:"updateAssetVersion"`
}

func (p *clientImpl) GetAsset(ctx context.Context, id int) (*Asset, error) {
	if id <= 0 {
		return nil, errors.New("asset id is required")
	}

	query := `query asset($id: ID!) { asset(id: $id) { id name type moduleId module { id abbreviation } versions { id status languageCode language } } }`

	variables := map[string]interface{}{
		"id": id,
	}

	var res GetAssetResponse
	if err := p.Exec(ctx, query, &res, variables); err != nil {
		return nil, err
	}

	return &res.Asset, nil
}

func (p *clientImpl) GetAssets(ctx context.Context, params AssetParams) ([]Asset, error) {
	query := `query assets($params: AssetParams) { assets(params: $params) { id name type moduleId module { id abbreviation } versions { id status languageCode language } } }`

	variables := map[string]interface{}{
		"params": map[string]interface{}{},
	}

	if params.ModuleID == 0 && params.ExternalModuleID == "" {
		return nil, errors.New("module id or external id is required")
	}

	if params.ModuleID > 0 {
		variables["params"].(map[string]interface{})["moduleId"] = params.ModuleID
	}

	if params.ExternalModuleID != "" {
		variables["params"].(map[string]interface{})["externalModuleId"] = params.ExternalModuleID
	}

	if params.Type != "" {
		variables["params"].(map[string]interface{})["type"] = params.Type
	}

	var res GetAssetsResponse
	if err := p.Exec(ctx, query, &res, variables); err != nil {
		return nil, err
	}

	return res.Assets, nil
}

func (p *clientImpl) CreateAsset(ctx context.Context, asset *Asset) error {
	if asset == nil {
		return errors.New("asset is nil")
	}

	query := `mutation createAsset($input: AssetInput!) { createAsset(input: $input) { id name type moduleId module { id abbreviation } } }`

	variables := map[string]interface{}{
		"input": map[string]interface{}{},
	}

	if asset.ModuleID > 0 {
		variables["input"].(map[string]interface{})["moduleId"] = asset.ModuleID
	}

	if asset.Name != "" {
		variables["input"].(map[string]interface{})["name"] = asset.Name
	}

	if asset.Type != "" {
		variables["input"].(map[string]interface{})["type"] = asset.Type
	}

	var res CreateAssetResponse
	if err := p.Exec(ctx, query, &res, variables); err != nil {
		return err
	}

	*asset = res.Asset
	return nil
}

func (p *clientImpl) CreateAssetVersion(ctx context.Context, assetVersion *AssetVersion) error {
	if assetVersion == nil {
		return errors.New("asset version is nil")
	}

	query := `mutation createAssetVersion($input: AssetVersionInput!) { createAssetVersion(input: $input) { id fileLink assetId asset { id name type moduleId module { id abbreviation } } } }`

	variables := map[string]interface{}{
		"input": map[string]interface{}{},
	}

	if assetVersion.AssetID > 0 {
		variables["input"].(map[string]interface{})["assetId"] = assetVersion.AssetID
	}

	if assetVersion.Status != "" {
		variables["input"].(map[string]interface{})["status"] = assetVersion.Status
	}

	if assetVersion.LanguageCode != "" {
		variables["input"].(map[string]interface{})["languageCode"] = assetVersion.LanguageCode
	}

	var assetVersionResponse CreateAssetVersionResponse
	if err := p.ExecWithFile(ctx, query, &assetVersionResponse, variables, assetVersion.LocalFilePath, "filePath"); err != nil {
		return err
	}

	*assetVersion = assetVersionResponse.AssetVersion
	return nil
}

func (p *clientImpl) UpdateAssetVersion(ctx context.Context, assetVersion *AssetVersion) error {
	if assetVersion == nil {
		return errors.New("asset version is nil")
	}

	query := `mutation updateAssetVersion($input: AssetVersionInput!) { updateAssetVersion(input: $input) { id status languageCode assetId asset { id name type moduleId module { id abbreviation } } } }`

	variables := map[string]interface{}{
		"input": map[string]interface{}{},
	}

	if assetVersion.ID > 0 {
		variables["input"].(map[string]interface{})["id"] = assetVersion.ID
	} else {
		return errors.New("asset version id is required")
	}

	if assetVersion.Status != "" {
		variables["input"].(map[string]interface{})["status"] = assetVersion.Status
	}

	if assetVersion.LanguageCode != "" {
		return errors.New("language code cannot be updated")
	}

	var assetVersionResponse UpdateAssetVersionResponse
	if err := p.Exec(ctx, query, &assetVersionResponse, variables); err != nil {
		return err
	}

	*assetVersion = assetVersionResponse.AssetVersion
	return nil
}
