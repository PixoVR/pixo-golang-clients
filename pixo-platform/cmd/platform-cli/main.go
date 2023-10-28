package main

import (
	"github.com/PixoVR/pixo-golang-clients/pixo-platform/cmd/platform-cli/parser"
	graphql_api "github.com/PixoVR/pixo-golang-clients/pixo-platform/graphql-api"
	"github.com/rs/zerolog/log"
	"github.com/spf13/pflag"
	"os"
)

var (
	apiClient       *graphql_api.GraphQLAPIClient
	iniFilePath     string
	semanticVersion string
	image           string
	moduleID        int
)

func init() {
	apiClient = graphql_api.NewClient(os.Getenv("SECRET_KEY"), "")
	pflag.StringVarP(&semanticVersion, "semantic-version", "v", "", "Multiplayer server version semantic version (e.g: 1.00.00)")
	pflag.StringVarP(&image, "image", "i", "", "Multiplayer server version image")
	pflag.StringVarP(&iniFilePath, "ini", "c", parser.DefaultConfigFilepath, "Path to ini config file")
	pflag.IntVarP(&moduleID, "module-id", "m", 0, "Multiplayer server version module id")
	pflag.Parse()
}

func main() {

	if moduleID == 0 {
		log.Fatal().Msg("Module id is required using -m <module-id> or --module-id <module-id>")
	}

	if image == "" {
		log.Fatal().Msg("Image is required using -i <image> or --image <image>")
	}

	iniParser, err := parser.NewIniParser(&iniFilePath)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to create ini parser")
		return
	}

	if semanticVersion == "" {
		semanticVersion, err = iniParser.ParseSemanticVersion()
		if err != nil {
			log.Fatal().Err(err).Msgf("No semantic version given and failed to parse server version from ini file %s", iniFilePath)
		}
	}

	if err = apiClient.CreateMultiplayerServerVersion(moduleID, image, semanticVersion); err != nil {
		log.Fatal().Err(err).Msg("Failed to update multiplayer server version")
	}

	log.Info().Msg("Successfully deployed multiplayer server version")
}
