package abstract_client

import (
	"fmt"
	"github.com/gorilla/websocket"
	"net/http"
	"time"
)

func (a *AbstractServiceClient) DialWebsocket(endpoint string) (*websocket.Conn, *http.Response, error) {

	header := http.Header{}
	if a.token != "" {
		header.Add("Authorization", "Bearer "+a.token)
	}

	url := a.GetURLWithPath(endpoint, "ws")
	conn, response, err := websocket.DefaultDialer.Dial(url, header)
	if err != nil {
		return nil, nil, fmt.Errorf("%w - %s", err, response.Status)
	}

	if err = conn.SetReadDeadline(time.Now().Add(a.timeoutDuration())); err != nil {
		return nil, nil, err
	}

	a.websocketConn = conn
	return conn, response, nil
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
