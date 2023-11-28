package graphql_api

import (
	"errors"
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
		log.Debug().Err(err).Msgf("error deploying multiplayer server version: %s", res)
		return err
	}

	log.Debug().Msgf("created multiplayer server version")
	return nil
}

func (g *GraphQLAPIClient) GetMultiplayerServerVersions(params MultiplayerServerVersionQueryParams) ([]*MultiplayerServerVersion, error) {

	configs, err := g.GetMultiplayerServerConfigs(MultiplayerServerConfigParams{
		ModuleID:      params.ModuleID,
		ServerVersion: params.SemanticVersion,
	})
	if err != nil {
		return nil, err
	}

	if len(configs) == 0 {
		return nil, errors.New("no multiplayer server configurations found")
	}

	return configs[0].ServerVersions, nil
}
