package matchmaker

import "net"

type Matchmaker interface {
	Connect(req MatchRequest) (*net.UDPAddr, error)
}
