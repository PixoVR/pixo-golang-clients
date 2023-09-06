package main

import (
	primary_api "github.com/PixoVR/pixo-golang-clients/pixo-platform/primary-api"
	"github.com/rs/zerolog/log"
	"net/http"
	"os"
	"strconv"
)

func main() {
	if len(os.Args) != 4 {
		log.Error().Msg("Invalid number of arguments. Expected 3 arguments: versionID, imageRegistry, status")
		return
	}

	versionIDString := os.Args[1]
	versionID, err := strconv.Atoi(versionIDString)
	if err != nil {
		log.Error().Err(err).Msg("Failed to parse versionID as int")
		return
	}

	imageRegistry := os.Args[2]
	status := os.Args[3]

	primaryClient := primary_api.NewClient(os.Getenv("SECRET_KEY"), "")
	res, err := primaryClient.UpdateMultiplayerServerVersion(versionID, imageRegistry, status)
	if err != nil {
		log.Error().Err(err).Msg("Failed to update multiplayer server version")
	}

	if res.StatusCode() != http.StatusOK {
		log.Error().Err(err).Msg("Failed to update multiplayer server version")
	}

	log.Info().Msg("Successfully updated multiplayer server version")
}
