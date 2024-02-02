package matchmaker

import (
	abstract_client "github.com/PixoVR/pixo-golang-clients/pixo-platform/abstract-client"
	"net"
)

type MockMatchmaker struct {
	*abstract_client.AbstractServiceClient
	NumCalledFindMatch int
}

func NewMockMatchmaker() *MockMatchmaker {
	return &MockMatchmaker{
		AbstractServiceClient: abstract_client.NewClient(abstract_client.AbstractConfig{}),
	}
}

func (m *MockMatchmaker) FindMatch(request MatchRequest) (*net.UDPAddr, error) {
	m.NumCalledFindMatch++
	return &net.UDPAddr{
		IP:   net.ParseIP(Localhost),
		Port: DefaultGameserverPort,
	}, nil
}

func (m *MockMatchmaker) Login(username, password string) error {
	return nil
}

func (m *MockMatchmaker) ActiveUserID() int {
	return 1
}

func (m *MockMatchmaker) DialGameserver(addr *net.UDPAddr) error {
	return nil
}

func (m *MockMatchmaker) CloseGameserverConnection() error {
	return nil
}

func (m *MockMatchmaker) SendAndReceiveMessage(message []byte) ([]byte, error) {
	return []byte("hello world"), nil
}
