package platform

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"time"

	"github.com/rs/zerolog/log"
)

type GitConfig struct {
	Provider string `json:"provider,omitempty"`
	OrgName  string `json:"orgName,omitempty"`
	RepoName string `json:"repoName,omitempty"`
}

type Language struct {
	Language     string `json:"language,omitempty"`
	LanguageCode string `json:"languageCode,omitempty"`
	DisplayName  string `json:"displayName,omitempty"`
}

type Module struct {
	ID           int    `json:"id,omitempty"`
	Name         string `json:"name,omitempty"`
	Abbreviation string `json:"abbreviation,omitempty"`
	Description  string `json:"description,omitempty"`
	ImageLink    string `json:"imageLink,omitempty"`
	ImagePath    string `json:"imagePath,omitempty"`
	PDFLink      string `json:"pdfLink,omitempty"`
	PDFPath      string `json:"pdfPath,omitempty"`
	ShortDesc    string `json:"shortDesc,omitempty"`
	LongDesc     string `json:"longDesc,omitempty"`
	Industry     string `json:"industry,omitempty"`
	Developer    string `json:"developer,omitempty"`
	Details      string `json:"details,omitempty"`
	Categories   string `json:"categories,omitempty"`
	Status       string `json:"status,omitempty"`
	ExternalID   string `json:"externalId,omitempty"`

	IsAvailable           bool `json:"isAvailable,omitempty"`
	IsPublic              bool `json:"isPublic,omitempty"`
	IsDemo                bool `json:"isDemo,omitempty"`
	IsMultiplayer         bool `json:"isMultiplayer,omitempty"`
	IsAuthenticatedLaunch bool `json:"isAuthenticatedLaunch,omitempty"`
	PassingScoreEnabled   bool `json:"passingScoreEnabled,omitempty"`

	ModulePlayerID int           `json:"modulePlayerId,omitempty"`
	ModulePlayer   *ModulePlayer `json:"modulePlayer,omitempty"`

	DistributorID int  `json:"distributorId,omitempty"`
	Distributor   *Org `json:"distributor,omitempty"`

	GitConfigID int       `json:"gitConfigId,omitempty"`
	GitConfig   GitConfig `json:"gitConfig,omitempty"`

	AvailableLanguages []Language      `json:"availableLanguages,omitempty"`
	Versions           []ModuleVersion `json:"versions,omitempty"`

	CreatedAt time.Time `json:"createdAt,omitempty"`
	UpdatedAt time.Time `json:"updatedAt,omitempty"`
}

type ModuleVersion struct {
	ID              int               `json:"id,omitempty"`
	ModuleID        int               `json:"moduleId,omitempty"`
	Module          Module            `json:"module,omitempty"`
	LifecycleID     int               `json:"lifecycleId,omitempty"`
	Lifecycle       *VersionLifecycle `json:"lifecycle,omitempty"`
	FileLink        string            `json:"fileLink,omitempty"`
	FilePath        string            `json:"filePath,omitempty"`
	FileSize        int               `json:"fileSize,omitempty"`
	SemanticVersion string            `json:"version,omitempty"`
	Notes           string            `json:"notes,omitempty"`
	Package         string            `json:"package,omitempty"`
	Public          bool              `json:"public,omitempty"`
	UploadStatus    string            `json:"uploadStatus,omitempty"`
	ExternalID      string            `json:"externalId,omitempty"`
	LocalFilePath   string            `json:"-"`
	ControlIds      []int             `json:"controlIds,omitempty"`
	PlatformIds     []int             `json:"platformIds,omitempty"`
	Platforms       []Platform        `json:"platforms,omitempty"`
	CreatedAt       time.Time         `json:"createdAt,omitempty"`
	UpdatedAt       time.Time         `json:"updatedAt,omitempty"`
}

// ModuleParams are the filters accepted by the modules query. Leaving a filter
// empty leaves that dimension unfiltered.
type ModuleParams struct {
	LifecycleIds  []int    `json:"lifecycleIds,omitempty"`
	Statuses      []string `json:"statuses,omitempty"`
	IsPublic      *bool    `json:"isPublic,omitempty"`
	DistributorID *int     `json:"distributorId,omitempty"`
}

// moduleFields is the selection every module read shares. Versions are left out
// so that listing modules stays a single cheap query - use GetModule or
// GetModuleVersions when versions are needed.
const moduleFields = `
	id
	name
	abbreviation
	description
	externalId
	imageLink
	imagePath
	pdfLink
	pdfPath
	shortDesc
	longDesc
	industry
	developer
	details
	categories
	status
	isAvailable
	isPublic
	isDemo
	isMultiplayer
	isAuthenticatedLaunch
	passingScoreEnabled
	modulePlayerId
	modulePlayer { id name description launchProtocol distributorId }
	distributorId
	distributor { id name type logoLink logoPath hubLogoLink }
	gitConfigId
	gitConfig { provider orgName repoName }
	availableLanguages { language languageCode displayName }
	createdAt
	updatedAt
`

type GetModulesResponse struct {
	Modules []Module `json:"modules"`
}

type GetModuleResponse struct {
	Module Module `json:"module"`
}

type CreateModuleResponse struct {
	Module Module `json:"createModule"`
}

type CreateModuleVersionResponse struct {
	ModuleVersion ModuleVersion `json:"createModuleVersion"`
}

func (p *clientImpl) GetModules(ctx context.Context, params ...ModuleParams) ([]Module, error) {
	query := fmt.Sprintf(`query modules($params: ModuleFilterParams) { modules(params: $params) { %s } }`, moduleFields)

	variables := map[string]interface{}{}
	if len(params) > 0 {
		variables["params"] = params[0]
	}

	var res GetModulesResponse
	if err := p.Exec(ctx, query, &res, variables); err != nil {
		return nil, err
	}

	return res.Modules, nil
}

// GetModule retrieves a single module along with its versions and the platforms
// each version supports.
func (p *clientImpl) GetModule(ctx context.Context, id int) (*Module, error) {
	if id == 0 {
		return nil, errors.New("module id is required")
	}

	query := fmt.Sprintf(
		`query module($id: ID!) { module(id: $id) { %s versions { %s } } }`,
		moduleFields,
		moduleVersionFields,
	)

	variables := map[string]interface{}{
		"id": id,
	}

	var res GetModuleResponse
	if err := p.Exec(ctx, query, &res, variables); err != nil {
		return nil, err
	}

	return &res.Module, nil
}

func (p *clientImpl) CreateModuleVersion(ctx context.Context, input ModuleVersion) (*ModuleVersion, error) {
	query := `mutation createModuleVersion($input: ModuleVersionInput!) { createModuleVersion(input: $input) { id moduleId module { abbreviation } version package lifecycleId lifecycle { id name } fileLink } }`

	if input.LocalFilePath == "" {
		return nil, errors.New("file path must be provided")
	}

	inputVariables := map[string]interface{}{
		"moduleId":    input.ModuleID,
		"version":     input.SemanticVersion,
		"notes":       input.Notes,
		"package":     input.Package,
		"platformIds": input.PlatformIds,
		"controlIds":  input.ControlIds,
	}

	if input.LifecycleID != 0 {
		inputVariables["lifecycleId"] = input.LifecycleID
	}

	variables := map[string]interface{}{
		"input": inputVariables,
	}

	if input.LocalFilePath == "" {
		var res CreateModuleVersionResponse
		if err := p.Exec(ctx, query, &res, variables); err != nil {
			return nil, err
		}

		return &res.ModuleVersion, nil
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
	defer func() { _ = file.Close() }()

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

	p.SetHeader("Content-Type", writer.FormDataContentType())

	res, err := p.Post(context.TODO(), "query", payload.Bytes())
	if err != nil {
		log.Error().Err(err).Msg("error creating multiplayer server version")
		return nil, err
	}

	resBody, _ := io.ReadAll(res.Body)

	if res.StatusCode > 299 {
		return nil, fmt.Errorf("error creating multiplayer server version: %s", string(resBody))
	}

	var gqlRes struct {
		Data   CreateModuleVersionResponse `json:"data"`
		Errors []struct {
			Message string `json:"message"`
		}
	}

	if err = json.Unmarshal(resBody, &gqlRes); err != nil {
		return nil, err
	}

	if len(gqlRes.Errors) > 0 {
		return nil, errors.New(gqlRes.Errors[0].Message)
	}

	return &gqlRes.Data.ModuleVersion, nil
}

func (p *clientImpl) GetModulesForUser(ctx context.Context, userID int) ([]Module, error) {
	query := fmt.Sprintf(`query user($id: ID!){ user(id: $id) { modules { %s } } }`, moduleFields)

	variables := map[string]interface{}{
		"id": userID,
	}

	var userResponse GetUserResponse
	if err := p.Exec(ctx, query, &userResponse, variables); err != nil {
		return nil, err
	}

	return userResponse.User.Modules, nil
}
