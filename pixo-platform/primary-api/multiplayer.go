package primary_api

import (
	"encoding/json"
	"github.com/go-resty/resty/v2"
	"github.com/rs/zerolog/log"
)

func (p *PrimaryAPIClient) UpdateMultiplayerServerVersion(imageRegistry, status string) (*resty.Response, error) {
	multiplayerPatch := MultiplayerServerVersion{
		Status:        status,
		ImageRegistry: imageRegistry,
	}
	body, err := json.Marshal(multiplayerPatch)
	if err != nil {
		log.Debug().Err(err).Msg("Failed to marshal multiplayer server version")
		return nil, err
	}

	return p.Patch("api/external/multiplayer-server-versions/1", body)
}
