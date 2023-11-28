package matchmaker

import "net"

type Matchmaker interface {
	FindMatch(req MatchRequest) (*net.UDPAddr, error)
	DialGameserver(addr *net.UDPAddr) error
	CloseGameserverConnection() error
	SendAndReceiveMessage(message []byte) ([]byte, error)
}
