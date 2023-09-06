package primary_api

import (
	"encoding/json"
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/rs/zerolog/log"
)

func (p *PrimaryAPIClient) DeployMultiplayerServerVersion(moduleID int, imageRegistry string) (*resty.Response, error) {
	multiplayerServerVersion := MultiplayerServerVersion{
		ModuleID:      moduleID,
		Status:        "enabled",
		ImageRegistry: imageRegistry,
	}

	body, err := json.Marshal(multiplayerServerVersion)
	if err != nil {
		log.Debug().Err(err).Msg("Failed to marshal multiplayer server version")
		return nil, err
	}

	path := "api/multiplayer-server-version"

	return p.Post(path, body)
}

func (p *PrimaryAPIClient) UpdateMultiplayerServerVersion(versionID int, imageRegistry string) (*resty.Response, error) {
	multiplayerPatch := MultiplayerServerVersion{
		Status:        "enabled",
		ImageRegistry: imageRegistry,
	}
	body, err := json.Marshal(multiplayerPatch)
	if err != nil {
		log.Debug().Err(err).Msg("Failed to marshal multiplayer server version")
		return nil, err
	}

	path := fmt.Sprintf("api/external/multiplayer-server-versions/%d", versionID)

	return p.Patch(path, body)
}
