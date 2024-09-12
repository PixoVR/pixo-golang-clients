package abstract_client

import (
	"fmt"
	"github.com/gorilla/websocket"
	"net/http"
	"time"
)

// DialWebsocket creates a new websocket connection to the given endpoint.
func (a *AbstractServiceClient) DialWebsocket(endpoint string) (*websocket.Conn, *http.Response, error) {
	url := a.GetURLWithPath(endpoint, "ws")
	conn, response, err := websocket.DefaultDialer.Dial(url, a.GetAuthHeader())
	if err != nil {
		return nil, nil, fmt.Errorf("%w - %s", err, response.Status)
	}

	if err = conn.SetReadDeadline(time.Now().Add(a.timeoutDuration())); err != nil {
		return nil, nil, err
	}

	a.websocketConn = conn
	return conn, response, nil
}

// WriteToWebsocket writes a message to the websocket connection.
func (a *AbstractServiceClient) WriteToWebsocket(message []byte) error {
	return a.websocketConn.WriteMessage(websocket.TextMessage, message)
}

// ReadFromWebsocket reads a message from the websocket connection.
func (a *AbstractServiceClient) ReadFromWebsocket() (int, []byte, error) {
	return a.websocketConn.ReadMessage()
}

// CloseWebsocketConnection closes the websocket connection.
func (a *AbstractServiceClient) CloseWebsocketConnection() error {
	return a.websocketConn.Close()
}

// SetWebsocketTimeout sets the timeout for the websocket connection. Defaults to 30 seconds.
func (a *AbstractServiceClient) timeoutDuration() time.Duration {
	if a.timeoutSeconds == 0 {
		return 30 * time.Second
	}

	return time.Duration(a.timeoutSeconds) * time.Second
}
