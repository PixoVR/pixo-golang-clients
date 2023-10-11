package matchmaker

import "net"

type Matchmaker interface {
	Connect(moduleID, orgID int) (*net.UDPAddr, error)
}
