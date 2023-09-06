package main

import (
	primary_api "github.com/PixoVR/pixo-golang-clients/pixo-platform/primary-api"
	"github.com/rs/zerolog/log"
	"net/http"
	"os"
	"strconv"
)

func main() {
	if len(os.Args) != 3 {
		log.Error().Msg("Invalid number of arguments. Expected 2 arguments: moduleID, imageRegistry")
		return
	}

	moduleID, err := strconv.Atoi(os.Args[1])
	if err != nil {
		log.Error().Err(err).Msg("Failed to parse moduleID as int")
		return
	}

	imageRegistry := os.Args[2]

	primaryClient := primary_api.NewClient(os.Getenv("SECRET_KEY"), "")
	res, err := primaryClient.DeployMultiplayerServerVersion(moduleID, imageRegistry)
	if err != nil {
		log.Error().Err(err).Msg("Failed to update multiplayer server version")
	}

	if res.StatusCode() != http.StatusOK {
		log.Error().Err(err).Msg("Failed to update multiplayer server version")
	}

	log.Info().Msg("Successfully updated multiplayer server version")
}
