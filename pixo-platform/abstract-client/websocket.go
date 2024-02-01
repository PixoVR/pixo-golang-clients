package abstract_client

import (
	"errors"
	"github.com/gorilla/websocket"
	"net/http"
	"time"
)

func (a *AbstractServiceClient) ConnectToWebsocket() (*http.Response, error) {

	httpHeader := http.Header{}
	if a.token != "" {
		httpHeader.Add("Authorization", "Bearer "+a.token)
	}

	conn, httpResponse, err := websocket.DefaultDialer.Dial(a.url, httpHeader)
	if err != nil {
		return nil, err
	}

	if err = conn.SetReadDeadline(time.Now().Add(a.timeoutDuration())); err != nil {
		return nil, err
	}

	a.conn = conn

	return httpResponse, nil
}

func (a *AbstractServiceClient) SendMessageToWebsocket(message []byte) error {

	if err := a.conn.WriteMessage(websocket.TextMessage, message); err != nil {
		return err
	}

	return nil
}

func (a *AbstractServiceClient) ReadFromWebsocket() ([]byte, error) {

	_, msg, err := a.conn.ReadMessage()
	if err != nil {
		return nil, err
	}

	if len(msg) == 0 {
		return nil, errors.New("empty response data")
	}

	return msg, nil
}

func (a *AbstractServiceClient) CloseWebsocketConnection() error {

	if err := a.conn.Close(); err != nil {
		return err
	}

	return nil
}

func (a *AbstractServiceClient) timeoutDuration() time.Duration {
	if a.timeoutSeconds == 0 {
		return 30 * time.Second
	}

	return time.Duration(a.timeoutSeconds) * time.Second
}
