package graphql_api

import (
	"context"
	platform "github.com/PixoVR/pixo-golang-clients/pixo-platform/primary-api"
	"github.com/rs/zerolog/log"
)

func (g *GraphQLAPIClient) DeployMultiplayerServerVersion(moduleID int, image, semanticVersion string) error {

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

	if _, err := g.gqlClient.ExecRaw(context.Background(), query, variables); err != nil {
		log.Debug().Msgf("error deploying multiplayer server version: %v", err)
		return err
	}

	log.Debug().Msgf("created multiplayer server version")
	return nil
}

func (g *GraphQLAPIClient) GetMultiplayerServerVersions() ([]*platform.MultiplayerServerVersion, error) {
	var res MultiplayerServerVersionQuery

	if err := g.gqlClient.Query(context.Background(), &res, nil); err != nil {
		return nil, err
	}

	return res.MultiplayerServerVersions, nil
}
