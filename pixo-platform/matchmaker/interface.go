package matchmaker

import (
	abstract_client "github.com/PixoVR/pixo-golang-clients/pixo-platform/abstract-client"
	"net"
)

type Matchmaker interface {
	abstract_client.AbstractClient
	FindMatch(req MatchRequest) (*net.UDPAddr, error)
	DialGameserver(addr *net.UDPAddr) error
	CloseGameserverConnection() error
	SendAndReceiveMessage(message []byte) ([]byte, error)
}
