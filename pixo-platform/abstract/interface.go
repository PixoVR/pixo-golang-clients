package abstract

import (
	"context"
	"github.com/gorilla/websocket"
	"net/http"
)

type AbstractClient interface {
	Path() string
	GetIPAddress() (string, error)
	Login(username, password string) error
	SetAPIKey(key string)
	SetToken(key string)
	GetToken() string
	GetURL(protocol ...string) string
	IsAuthenticated() bool

	Get(ctx context.Context, path string) (*http.Response, error)
	Post(ctx context.Context, path string, body []byte) (*http.Response, error)
	Put(ctx context.Context, path string, body []byte) (*http.Response, error)
	Patch(ctx context.Context, path string, body []byte) (*http.Response, error)
	Delete(ctx context.Context, path string) (*http.Response, error)

	DialWebsocket(endpoint string) (*websocket.Conn, *http.Response, error)
	WriteToWebsocket(message []byte) error
	ReadFromWebsocket() (int, []byte, error)
	CloseWebsocketConnection() error
}
