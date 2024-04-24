package matchmaker

import (
	abstractClient "github.com/PixoVR/pixo-golang-clients/pixo-platform/abstract-client"
	graphql_api "github.com/PixoVR/pixo-golang-clients/pixo-platform/graphql-api"
	"github.com/PixoVR/pixo-golang-clients/pixo-platform/urlfinder"
	"net"
)

type MultiplayerMatchmaker struct {
	*abstractClient.AbstractServiceClient
	gameserverAddress    *net.UDPAddr
	gameserverConnection *net.UDPConn
}

func NewMatchmakerWithBasicAuth(username, password string, config urlfinder.ClientConfig, timeoutSeconds ...int) (*MultiplayerMatchmaker, error) {
	platformClient, err := graphql_api.NewClientWithBasicAuth(username, password, config)
	if err != nil {
		return nil, err
	}

	config.Token = platformClient.GetToken()

	return NewMatchmaker(config, timeoutSeconds...), nil
}

func NewMatchmaker(config urlfinder.ClientConfig, timeoutSeconds ...int) *MultiplayerMatchmaker {

	if len(timeoutSeconds) == 0 {
		timeoutSeconds = []int{60}
	}

	serviceConfig := newServiceConfig(config.Lifecycle, config.Region)
	url := getURL(serviceConfig.FormatURL())
	abstractConfig := abstractClient.AbstractConfig{
		URL:            url,
		Token:          config.Token,
		TimeoutSeconds: timeoutSeconds[0],
	}

	return &MultiplayerMatchmaker{
		AbstractServiceClient: abstractClient.NewClient(abstractConfig),
	}
}

func (m *MultiplayerMatchmaker) Login(username, password string) error {
	return nil
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

	return host
}
