package abstract_client

import (
	"net/http"
)

type AbstractClient interface {
	Login(username, password string) error
	SetAPIKey(key string)
	SetToken(key string)
	GetToken() string
	GetURL() string
	IsAuthenticated() bool
	ActiveUserID() int

	DialWebsocket(endpoint string) (*http.Response, error)
	WriteToWebsocket(message []byte) error
	ReadFromWebsocket() (int, []byte, error)
	CloseWebsocketConnection() error
}
