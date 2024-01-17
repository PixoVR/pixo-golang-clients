package matchmaker

import (
	"fmt"
	abstractClient "github.com/PixoVR/pixo-golang-clients/pixo-platform/abstract-client"
	primary_api "github.com/PixoVR/pixo-golang-clients/pixo-platform/primary-api"
	"github.com/PixoVR/pixo-golang-clients/pixo-platform/urlfinder"
	"net"
)

type MultiplayerMatchmaker struct {
	abstractClient.PixoAbstractAPIClient
	gameserverAddress    *net.UDPAddr
	gameserverConnection *net.UDPConn
}

func NewMatchmakerWithBasicAuth(username, password, lifecycle, region string, timeoutSeconds ...int) (*MultiplayerMatchmaker, error) {
	primaryClient, err := primary_api.NewClientWithBasicAuth(username, password, lifecycle, region)
	if err != nil {
		return nil, err
	}

	return NewMatchmaker(lifecycle, region, primaryClient.GetToken(), timeoutSeconds...), nil
}

func NewMatchmaker(lifecycle, region, token string, timeoutSeconds ...int) *MultiplayerMatchmaker {

	if len(timeoutSeconds) == 0 {
		timeoutSeconds = []int{60}
	}

	config := newServiceConfig(lifecycle, region)
	url := getURL(config.FormatURL())

	return &MultiplayerMatchmaker{
		PixoAbstractAPIClient: *abstractClient.NewClient(token, url, timeoutSeconds[0]),
	}
}

func newServiceConfig(lifecycle, region string) urlfinder.ServiceConfig {
	return urlfinder.ServiceConfig{
		Service:   "match",
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
