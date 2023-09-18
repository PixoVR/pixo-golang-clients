package main

import (
	"github.com/PixoVR/pixo-golang-clients/pixo-platform/cmd/platform-cli/parser"
	primary_api "github.com/PixoVR/pixo-golang-clients/pixo-platform/primary-api"
	"github.com/go-resty/resty/v2"
	"github.com/rs/zerolog/log"
	"net/http"
	"os"
	"strconv"
)

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

	res, err := deployMultiplayerServerVersion(id, image, semanticVersion)
	if err != nil {
		log.Error().Err(err).Msg("Failed to update multiplayer server version")
		return
	}

	if res.StatusCode() != http.StatusOK {
		log.Error().Err(err).Msg("Failed to update multiplayer server version")
		return
	}

	log.Info().Msg("Successfully deployed multiplayer server version")
}

func deployMultiplayerServerVersion(id int, image, semanticVersion string) (*resty.Response, error) {
	primaryClient := primary_api.NewClientWithBasicAuth(os.Getenv("PIXO_USERNAME"), os.Getenv("PIXO_PASSWORD"), "")
	return primaryClient.DeployMultiplayerServerVersion(id, image, semanticVersion)
}
