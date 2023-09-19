package graphql_api

import platform "github.com/PixoVR/pixo-golang-clients/pixo-platform/primary-api"

type MultiplayerServerVersionQuery struct {
	MultiplayerServerVersions []*platform.MultiplayerServerVersion `graphql:"multiplayerServerVersions"`
}

type MultiplayerServerVersionInput struct {
	Input platform.MultiplayerServerVersion `graphql:"createMultiplayerServer($input: MultiplayerServerVersionInput!)"`
}
type Response struct {
	CreateMultiplayerServerVersion platform.MultiplayerServerVersion `graphql:"createMultiplayerServerVersion(input: $input)"`
}
