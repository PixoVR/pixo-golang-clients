package primary_api

import (
	"encoding/json"
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/rs/zerolog/log"
)

func (p *PrimaryAPIClient) DeployMultiplayerServerVersion(moduleID int, image, semanticVersion string) (*resty.Response, error) {
	multiplayerServerVersion := MultiplayerServerVersion{
		ModuleID:         moduleID,
		Status:           "enabled",
		ImageRegistry:    image,
		Engine:           "unreal",
		Version:          semanticVersion,
		MinClientVersion: semanticVersion,
		Filename:         "nonexistent-file.exe",
	}

	body, err := json.Marshal(multiplayerServerVersion)
	if err != nil {
		log.Debug().Err(err).Msg("Failed to marshal multiplayer server version")
		return nil, err
	}

	path := "api/multiplayer-server-version"

	return p.Post(path, body)
}

func (p *PrimaryAPIClient) UpdateMultiplayerServerVersion(versionID int, image string) (*resty.Response, error) {
	multiplayerPatch := MultiplayerServerVersion{
		Status:        "enabled",
		ImageRegistry: image,
	}

	body, err := json.Marshal(multiplayerPatch)
	if err != nil {
		log.Debug().Err(err).Msg("Failed to marshal multiplayer server version")
		return nil, err
	}

	path := fmt.Sprintf("api/external/multiplayer-server-versions/%d", versionID)

	return p.Patch(path, body)
}

func (p *PrimaryAPIClient) GetMultiplayerConfigurations() (*resty.Response, error) {
	path := "api/openmatch/configurations"

	p.AddHeader("x-openmatch-header", p.GetToken())

	return p.Get(path)
}
