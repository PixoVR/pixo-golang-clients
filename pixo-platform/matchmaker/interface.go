package matchmaker

import (
	abstract_client "github.com/PixoVR/pixo-golang-clients/pixo-platform/abstract-client"
	"github.com/gorilla/websocket"
	"net"
	"net/http"
)

type Matchmaker interface {
	abstract_client.AbstractClient
	FindMatch(req MatchRequest) (*net.UDPAddr, error)

	DialMatchmaker() (*websocket.Conn, *http.Response, error)
	SendRequest(conn *websocket.Conn, req MatchRequest) error
	ReadResponse(conn *websocket.Conn) (MatchResponse, error)
	CloseMatchmakerConnection(conn *websocket.Conn) error

	DialGameserver(addr *net.UDPAddr) error
	CloseGameserverConnection() error
	SendAndReceiveMessage(message []byte) ([]byte, error)
}
