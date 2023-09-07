package main

import (
	primary_api "github.com/PixoVR/pixo-golang-clients/pixo-platform/primary-api"
	"github.com/go-resty/resty/v2"
	"github.com/rs/zerolog/log"
	"net/http"
	"os"
	"strconv"
)

func main() {
	if len(os.Args) != 4 {
		log.Error().Msg("Invalid number of arguments. Expected 2 arguments: id, image, semantic version (e.g. 1.00.00)")
		return
	}

	id, err := strconv.Atoi(os.Args[1])
	if err != nil {
		log.Error().Err(err).Msg("Failed to parse id as int")
		return
	}

	image := os.Args[2]
	semanticVersion := os.Args[3]

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
