package legacy

import (
	"encoding/json"
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/rs/zerolog/log"
)

func (p *Client) DeployMultiplayerServerVersion(moduleID int, image, semanticVersion string) (*resty.Response, error) {
	multiplayerServerVersion := MultiplayerServerVersion{
		ModuleID:        moduleID,
		Status:          "enabled",
		ImageRegistry:   image,
		Engine:          "unreal",
		SemanticVersion: semanticVersion,
	}

	body, err := json.Marshal(multiplayerServerVersion)
	if err != nil {
		log.Debug().Err(err).Msg("Failed to marshal multiplayer server version")
		return nil, err
	}

	path := "api/multiplayer-server-version"

	return p.Post(path, body)
}

func (p *Client) UpdateMultiplayerServerVersion(versionID int, image string) (*resty.Response, error) {
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

func (p *Client) GetMatchmakingProfiles() ([]*GameProfileMetadata, error) {
	path := "api/openmatch/configurations"

	p.SetHeader("x-openmatch-header", p.GetToken())

	res, err := p.Get(path)
	if err != nil {
		log.Debug().Err(err).Msg("Failed to get multiplayer configurations")
		return nil, err
	}

	var profilesResponse GameProfileMetaDataResponse
	if err = json.Unmarshal(res.Body(), &profilesResponse); err != nil {
		log.Debug().Err(err).Msg("Failed to unmarshal multiplayer configurations")
		return nil, err
	}

	return profilesResponse.Profiles, nil
}
