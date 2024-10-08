package matchmaker

import (
	abstract_client "github.com/PixoVR/pixo-golang-clients/pixo-platform/abstract"
	"github.com/gorilla/websocket"
	"net"
	"net/http"
)

// Matchmaker is the interface for the matchmaker client
type Matchmaker interface {
	abstract_client.AbstractClient
	FindMatch(req MatchRequest) (*net.UDPAddr, error)

	DialMatchmaker() (*websocket.Conn, *http.Response, error)
	SendRequest(conn *websocket.Conn, req MatchRequest) error
	ReadResponse(conn *websocket.Conn) (MatchResponse, error)
	CloseMatchmakerConnection(conn *websocket.Conn) error

	DialGameserver(addr *net.UDPAddr) error
	CloseGameserverConnection() error
	SendMessageToGameserver(message []byte) error
	ReadMessageFromGameserver() ([]byte, error)
}
