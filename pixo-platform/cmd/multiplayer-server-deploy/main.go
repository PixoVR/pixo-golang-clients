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
	res, err := updateMultiplayerServerVersion()
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

func deployMultiplayerServerVersion() (*resty.Response, error) {
	if len(os.Args) != 3 {
		log.Error().Msg("Invalid number of arguments. Expected 2 arguments: moduleID, imageRegistry")
		return nil, nil
	}

	moduleID, err := strconv.Atoi(os.Args[1])
	if err != nil {
		log.Error().Err(err).Msg("Failed to parse moduleID as int")
		return nil, nil
	}

	imageRegistry := os.Args[2]

	primaryClient := primary_api.NewClientWithBasicAuth(os.Getenv("PIXO_USERNAME"), os.Getenv("PIXO_PASSWORD"), "")
	return primaryClient.DeployMultiplayerServerVersion(moduleID, imageRegistry)
}

func updateMultiplayerServerVersion() (*resty.Response, error) {
	if len(os.Args) != 3 {
		log.Error().Msg("Invalid number of arguments. Expected 2 arguments: versionID, imageRegistry")
		return nil, nil
	}

	versionID, err := strconv.Atoi(os.Args[1])
	if err != nil {
		log.Error().Err(err).Msg("Failed to parse versionID as int")
		return nil, nil
	}

	imageRegistry := os.Args[2]

	secretKeyClient := primary_api.NewClient(os.Getenv("SECRET_KEY"), "")
	return secretKeyClient.UpdateMultiplayerServerVersion(versionID, imageRegistry)
}
