package matchmaker

import (
	"net"
)

type MockMatchmaker struct{}

func (m *MockMatchmaker) Connect(moduleID, orgID int) (*net.UDPAddr, error) {
	return &net.UDPAddr{
		IP:   net.ParseIP(Localhost),
		Port: DefaultGameserverPort,
	}, nil
}
