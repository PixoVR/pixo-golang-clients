package matchmaker

import (
	abstractClient "github.com/PixoVR/pixo-golang-clients/pixo-platform/abstract-client"
	graphql_api "github.com/PixoVR/pixo-golang-clients/pixo-platform/graphql-api"
	"github.com/PixoVR/pixo-golang-clients/pixo-platform/urlfinder"
	"net"
)

type MultiplayerMatchmaker struct {
	*abstractClient.AbstractServiceClient

	platformClient graphql_api.PlatformClient

	gameserverAddress    *net.UDPAddr
	gameserverConnection *net.UDPConn
}

func NewClientWithBasicAuth(username, password string, config urlfinder.ClientConfig, timeoutSeconds ...int) (*MultiplayerMatchmaker, error) {
	platformClient, err := graphql_api.NewClientWithBasicAuth(username, password, config)
	if err != nil {
		return nil, err
	}

	config.Token = platformClient.GetToken()

	return NewClient(config, timeoutSeconds...), nil
}

func NewClient(config urlfinder.ClientConfig, timeoutSeconds ...int) *MultiplayerMatchmaker {

	if len(timeoutSeconds) == 0 {
		timeoutSeconds = []int{60}
	}

	serviceConfig := newServiceConfig(config.Lifecycle, config.Region)
	abstractConfig := abstractClient.AbstractConfig{
		ServiceConfig:  serviceConfig,
		Token:          config.Token,
		TimeoutSeconds: timeoutSeconds[0],
	}

	return &MultiplayerMatchmaker{
		AbstractServiceClient: abstractClient.NewClient(abstractConfig),
		platformClient:        graphql_api.NewClient(config),
	}
}

func (m *MultiplayerMatchmaker) Login(username, password string) error {
	return m.platformClient.Login(username, password)
}

func newServiceConfig(lifecycle, region string) urlfinder.ServiceConfig {
	return urlfinder.ServiceConfig{
		Service:   "matchmaking",
		Lifecycle: lifecycle,
		Region:    region,
		Port:      8080,
	}
}
