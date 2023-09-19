package main

import (
	"github.com/PixoVR/pixo-golang-clients/pixo-platform/cmd/platform-cli/parser"
	graphql_api "github.com/PixoVR/pixo-golang-clients/pixo-platform/graphql-api"
	"github.com/rs/zerolog/log"
	"os"
	"strconv"
)

var (
	apiClient *graphql_api.GraphQLAPIClient
)

func init() {
	apiClient = graphql_api.NewClient(os.Getenv("SECRET_KEY"), "")
}

func main() {
	if len(os.Args) < 3 {
		log.Error().Msg("Invalid number of arguments. Expected at least 2 arguments: id, image, and optionally an ini config file path (default: ./Config/DefaultGame.ini)")
		return
	}

	id, err := strconv.Atoi(os.Args[1])
	if err != nil {
		log.Error().Err(err).Msg("Failed to parse id as int")
		return
	}

	image := os.Args[2]

	var ini *string
	if len(os.Args) < 4 {
		log.Info().Msg("No ini config file path provided. Using default: ./Config/DefaultGame.ini")
	} else {
		ini = &os.Args[3]
	}

	iniParser, err := parser.NewIniParser(ini)
	if err != nil {
		log.Error().Err(err).Msg("Failed to create ini parser")
		return
	}

	semanticVersion, err := iniParser.ParseServerVersion()
	if err != nil {
		log.Error().Err(err).Msg("Failed to parse server version")
		return
	}

	if err = apiClient.DeployMultiplayerServerVersion(id, image, semanticVersion); err != nil {
		log.Error().Err(err).Msg("Failed to update multiplayer server version")
		return
	}

	log.Info().Msg("Successfully deployed multiplayer server version")
}
