package abstract_client

import (
	"github.com/go-resty/resty/v2"
	"github.com/gorilla/websocket"
	"net/http"
)

type AbstractClient interface {
	GetIPAddress() (string, error)
	Login(username, password string) error
	SetAPIKey(key string)
	SetToken(key string)
	GetToken() string
	GetURL() string
	IsAuthenticated() bool

	Get(path string) (*resty.Response, error)
	Post(path string, body []byte) (*resty.Response, error)
	Put(path string, body []byte) (*resty.Response, error)
	Patch(path string, body []byte) (*resty.Response, error)
	Delete(path string) (*resty.Response, error)

	DialWebsocket(endpoint string) (*websocket.Conn, *http.Response, error)
	WriteToWebsocket(message []byte) error
	ReadFromWebsocket() (int, []byte, error)
	CloseWebsocketConnection() error
}
