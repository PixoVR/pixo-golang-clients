package matchmaker

import (
	"encoding/json"
	"errors"
	"github.com/gorilla/websocket"
	"io"
	"net"
	"net/http"
	"strconv"
)

// FindMatch returns the UDP address of the game server returned from the matchmaking request
func (m *MultiplayerMatchmaker) FindMatch(req MatchRequest) (*net.UDPAddr, error) {
	if !req.IsValid() {
		return nil, errors.New("match request is invalid")
	}

	_, httpResponse, err := m.DialWebsocket(MatchmakingEndpoint)
	if err != nil {
		return nil, err
	}

	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(httpResponse.Body)

	message, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	if err = m.WriteToWebsocket(message); err != nil {
		return nil, err
	}

	_, response, err := m.ReadFromWebsocket()
	if err != nil {
		return nil, err
	}

	var matchResponse MatchResponse
	if err = json.Unmarshal(response, &matchResponse); err != nil {
		return nil, err
	}

	if !matchResponse.IsValid() {
		return nil, errors.New(matchResponse.Message)
	}

	port, err := strconv.Atoi(matchResponse.MatchDetails.Port)
	if err != nil {
		return nil, err
	}

	m.gameserverAddress = &net.UDPAddr{
		IP:   net.ParseIP(matchResponse.MatchDetails.IP),
		Port: port,
	}
	return m.gameserverAddress, nil
}

func (m *MultiplayerMatchmaker) DialMatchmaker() (*websocket.Conn, *http.Response, error) {
	return m.DialWebsocket(MatchmakingEndpoint)
}

func (m *MultiplayerMatchmaker) SendRequest(conn *websocket.Conn, req MatchRequest) error {
	message, err := json.Marshal(req)
	if err != nil {
		return err
	}

	return conn.WriteMessage(websocket.TextMessage, message)
}

func (m *MultiplayerMatchmaker) ReadResponse(conn *websocket.Conn) (MatchResponse, error) {
	_, response, err := conn.ReadMessage()
	if err != nil {
		return MatchResponse{}, err
	}

	var matchResponse MatchResponse
	if err = json.Unmarshal(response, &matchResponse); err != nil {
		return MatchResponse{}, err
	}

	return matchResponse, nil
}

func (m *MultiplayerMatchmaker) CloseMatchmakerConnection(conn *websocket.Conn) error {
	return conn.Close()
}
