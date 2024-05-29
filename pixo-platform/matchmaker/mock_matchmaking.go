package matchmaker

import (
	"encoding/json"
	"fmt"
	abstract_client "github.com/PixoVR/pixo-golang-clients/pixo-platform/abstract-client"
	"github.com/gorilla/websocket"
	"net"
	"net/http"
)

type MockMatchmaker struct {
	abstract_client.MockAbstractClient

	NumCalledDialMatchmaker int
	DialMatchmakerError     error

	NumCalledFindMatch int
	FindMatchError     error

	NumCalledReadResponse int
	ReadResponseError     error

	NumCalledDialGameserver int
	DialGameserverError     error

	NumCalledSendToGameserver int
	SendToGameserverError     error

	NumCalledReadFromGameserver int
	ReadFromGameserverError     error

	NumCalledCloseGameserver int
	CloseGameserverError     error

	response []byte
}

func NewMockMatchmaker() *MockMatchmaker {
	mockResponse := MatchResponse{
		Message: "Match found",
		MatchDetails: MatchDetails{
			IP:   Localhost,
			Port: fmt.Sprint(DefaultGameserverPort),
		},
	}

	mockResponseBytes, _ := json.Marshal(mockResponse)

	return &MockMatchmaker{
		MockAbstractClient: abstract_client.MockAbstractClient{
			Response: mockResponseBytes,
		},
	}
}

func (m *MockMatchmaker) FindMatch(request MatchRequest) (*net.UDPAddr, error) {
	m.NumCalledFindMatch++

	if m.FindMatchError != nil {
		return nil, m.FindMatchError
	}

	return &net.UDPAddr{
		IP:   net.ParseIP(Localhost),
		Port: DefaultGameserverPort,
	}, nil
}

func (m *MockMatchmaker) DialMatchmaker() (*websocket.Conn, *http.Response, error) {
	m.NumCalledDialMatchmaker++

	if m.DialMatchmakerError != nil {
		return nil, nil, m.DialMatchmakerError

	}
	return nil, nil, nil
}

func (m *MockMatchmaker) SendRequest(conn *websocket.Conn, req MatchRequest) error {
	reqBytes, _ := json.Marshal(req)

	if m.WriteToWebsocketError != nil {
		return m.WriteToWebsocketError
	}

	return m.WriteToWebsocket(reqBytes)
}

func (m *MockMatchmaker) ReadResponse(conn *websocket.Conn) (MatchResponse, error) {
	m.NumCalledReadResponse++

	if m.ReadResponseError != nil {
		return MatchResponse{}, m.ReadResponseError
	}

	response := MatchResponse{}
	if err := json.Unmarshal(m.response, &response); err != nil {
		return MatchResponse{}, err
	}

	return response, nil
}

func (m *MockMatchmaker) CloseMatchmakerConnection(conn *websocket.Conn) error {
	return m.CloseWebsocketConnection()
}

func (m *MockMatchmaker) DialGameserver(addr *net.UDPAddr) error {
	m.NumCalledDialGameserver++

	if m.DialGameserverError != nil {
		return m.DialGameserverError
	}

	return nil
}

func (m *MockMatchmaker) CloseGameserverConnection() error {
	m.NumCalledCloseGameserver++

	if m.CloseGameserverError != nil {
		return m.CloseGameserverError
	}

	return nil
}

func (m *MockMatchmaker) SendMessageToGameserver(message []byte) error {
	m.NumCalledSendToGameserver++

	if m.SendToGameserverError != nil {
		return m.SendToGameserverError
	}

	return nil
}

func (m *MockMatchmaker) ReadMessageFromGameserver() ([]byte, error) {
	m.NumCalledReadFromGameserver++

	if m.ReadFromGameserverError != nil {
		return nil, m.ReadFromGameserverError
	}

	return []byte("{\"IPAddress\":\"127.0.0.1\",\"Port\":7777}"), nil
}
