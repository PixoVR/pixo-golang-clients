package matchmaker

import (
	"fmt"
	abstractClient "github.com/PixoVR/pixo-golang-clients/pixo-platform/abstract-client"
)

func NewMatchmaker(url, token string) *PixoMatchmaker {

	if url == "" {
		url = getURL()
	}

	return &PixoMatchmaker{
		PixoAbstractAPIClient: *abstractClient.NewClient(token, url),
	}
}

func getURL() string {
	return fmt.Sprintf("%s/%s", DefaultMatchmakingURL, MatchmakingEndpoint)
}
