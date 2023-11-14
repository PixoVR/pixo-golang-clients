package graphql_api

import (
	"github.com/rs/zerolog/log"
)

func (g *GraphQLAPIClient) GetMultiplayerServerConfigs(params MultiplayerServerConfigParams) ([]*MultiplayerServerConfigQueryParams, error) {
	var query MultiplayerServerConfigQuery

	variables := map[string]interface{}{
		"params": params,
	}

	if err := g.Query(&query, variables); err != nil {
		return nil, err
	}

	return query.MultiplayerServerConfigs, nil
}

func (g *GraphQLAPIClient) CreateMultiplayerServerVersion(moduleID int, image, semanticVersion string) error {
	query := `mutation createMultiplayerServerVersion($input: MultiplayerServerVersionInput!) { createMultiplayerServerVersion(input: $input) { id imageRegistry semanticVersion module { name } } }`

	variables := map[string]interface{}{
		"input": map[string]interface{}{
			"moduleId":        moduleID,
			"imageRegistry":   image,
			"semanticVersion": semanticVersion,
			"status":          "enabled",
			"engine":          "unreal",
		},
	}

	if res, err := g.ExecRaw(query, variables); err != nil {
		log.Debug().Msgf("error deploying multiplayer server version: %s", res)
		return err
	}

	log.Debug().Msgf("created multiplayer server version")
	return nil
}
