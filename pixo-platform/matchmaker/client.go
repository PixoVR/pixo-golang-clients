package matchmaker

import (
	"fmt"
	abstractClient "github.com/PixoVR/pixo-golang-clients/pixo-platform/abstract-client"
	graphql_api "github.com/PixoVR/pixo-golang-clients/pixo-platform/graphql-api"
	"github.com/PixoVR/pixo-golang-clients/pixo-platform/urlfinder"
	"net"
)

type MultiplayerMatchmaker struct {
	abstractClient.PixoAbstractAPIClient
	gameserverAddress    *net.UDPAddr
	gameserverConnection *net.UDPConn
}

func NewMatchmakerWithBasicAuth(username, password, lifecycle, region string, timeoutSeconds ...int) (*MultiplayerMatchmaker, error) {
	config := urlfinder.ClientConfig{
		Lifecycle: lifecycle,
		Region:    region,
	}
	platformClient, err := graphql_api.NewClientWithBasicAuth(username, password, config)
	if err != nil {
		return nil, err
	}

	return NewMatchmaker(lifecycle, region, platformClient.GetToken(), timeoutSeconds...), nil
}

func NewMatchmaker(lifecycle, region, token string, timeoutSeconds ...int) *MultiplayerMatchmaker {

	if len(timeoutSeconds) == 0 {
		timeoutSeconds = []int{60}
	}

	config := newServiceConfig(lifecycle, region)
	url := getURL(config.FormatURL())
	abstractConfig := abstractClient.AbstractConfig{
		URL:            url,
		Token:          token,
		TimeoutSeconds: timeoutSeconds[0],
	}

	return &MultiplayerMatchmaker{
		PixoAbstractAPIClient: *abstractClient.NewClient(abstractConfig),
	}
}

func newServiceConfig(lifecycle, region string) urlfinder.ServiceConfig {
	return urlfinder.ServiceConfig{
		Service:   "matchmaking",
		Lifecycle: lifecycle,
		Region:    region,
		Port:      8080,
	}
}

func getURL(host string) string {
	if host == "" {
		host = DefaultMatchmakingURL
	}

	return fmt.Sprintf("%s/%s", host, MatchmakingEndpoint)
}
