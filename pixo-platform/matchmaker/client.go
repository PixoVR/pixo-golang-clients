package matchmaker

import (
	abstract "github.com/PixoVR/pixo-golang-clients/pixo-platform/abstract-client"
	"github.com/PixoVR/pixo-golang-clients/pixo-platform/platform"
	"github.com/PixoVR/pixo-golang-clients/pixo-platform/urlfinder"
	"net"
	"sync"
)

type MultiplayerMatchmaker struct {
	*abstract.AbstractServiceClient
	*sync.Mutex

	platformClient platform.Client

	gameserverAddress    *net.UDPAddr
	gameserverConnection *net.UDPConn
}

func NewClientWithBasicAuth(username, password string, config urlfinder.ClientConfig, timeoutSeconds ...int) (*MultiplayerMatchmaker, error) {
	platformClient, err := platform.NewClientWithBasicAuth(username, password, config)
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
	abstractConfig := abstract.AbstractConfig{
		ServiceConfig:  serviceConfig,
		Token:          config.Token,
		TimeoutSeconds: timeoutSeconds[0],
	}

	return &MultiplayerMatchmaker{
		AbstractServiceClient: abstract.NewClient(abstractConfig),
		platformClient:        platform.NewClient(config),
		Mutex:                 &sync.Mutex{},
	}
}

func (m *MultiplayerMatchmaker) Login(username, password string) error {
	m.Lock()
	defer m.Unlock()
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
