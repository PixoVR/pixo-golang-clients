package abstract_client

import (
	"github.com/gorilla/websocket"
	"net/http"
	"time"
)

func (a *AbstractServiceClient) DialWebsocket(endpoint string) (*http.Response, error) {

	httpHeader := http.Header{}
	if a.token != "" {
		httpHeader.Add("Authorization", "Bearer "+a.token)
	}

	conn, httpResponse, err := websocket.DefaultDialer.Dial(a.GetURLWithPath(endpoint), httpHeader)
	if err != nil {
		return nil, err
	}

	if err = conn.SetReadDeadline(time.Now().Add(a.timeoutDuration())); err != nil {
		return nil, err
	}

	a.websocketConn = conn
	return httpResponse, nil
}

func (a *AbstractServiceClient) WriteToWebsocket(message []byte) error {
	return a.websocketConn.WriteMessage(websocket.TextMessage, message)
}

func (a *AbstractServiceClient) ReadFromWebsocket() (int, []byte, error) {
	return a.websocketConn.ReadMessage()
}

func (a *AbstractServiceClient) CloseWebsocketConnection() error {
	return a.websocketConn.Close()
}

func (a *AbstractServiceClient) timeoutDuration() time.Duration {
	if a.timeoutSeconds == 0 {
		return 30 * time.Second
	}

	return time.Duration(a.timeoutSeconds) * time.Second
}
