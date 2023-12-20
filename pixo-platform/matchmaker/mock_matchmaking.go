package matchmaker

import (
	"net"
)

type MockMatchmaker struct{}

func (m *MockMatchmaker) FindMatch(request MatchRequest) (*net.UDPAddr, error) {
	return &net.UDPAddr{
		IP:   net.ParseIP(Localhost),
		Port: DefaultGameserverPort,
	}, nil
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
