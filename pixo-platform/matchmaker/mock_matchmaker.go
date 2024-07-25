package matchmaker

import (
	"encoding/json"
	"fmt"
	abstract_client "github.com/PixoVR/pixo-golang-clients/pixo-platform/abstract-client"
	"github.com/gorilla/websocket"
	"net"
	"net/http"
	"sync"
)

type MockMatchmaker struct {
	abstract_client.MockAbstractClient
	*sync.Mutex

	NumCalledFindMatch int
	FindMatchError     error

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
		Mutex: &sync.Mutex{},
		MockAbstractClient: abstract_client.MockAbstractClient{
			Response: mockResponseBytes,
		},
	}
}

func (m *MockMatchmaker) Reset() {
	m.Lock()
	defer m.Unlock()

	m.NumCalledFindMatch = 0
	m.FindMatchError = nil

	m.NumCalledDialGameserver = 0
	m.DialGameserverError = nil

	m.NumCalledSendToGameserver = 0
	m.SendToGameserverError = nil

	m.NumCalledReadFromGameserver = 0
	m.ReadFromGameserverError = nil

	m.NumCalledCloseGameserver = 0
	m.CloseGameserverError = nil

	m.response = nil
}

func (m *MockMatchmaker) FindMatch(request MatchRequest) (*net.UDPAddr, error) {
	m.Lock()
	defer m.Unlock()

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
	m.Lock()
	defer m.Unlock()
	return m.DialWebsocket("")
}

func (m *MockMatchmaker) SendRequest(conn *websocket.Conn, req MatchRequest) error {
	m.Lock()
	defer m.Unlock()

	reqBytes, _ := json.Marshal(req)
	return m.WriteToWebsocket(reqBytes)
}

func (m *MockMatchmaker) ReadResponse(conn *websocket.Conn) (MatchResponse, error) {
	m.Lock()
	defer m.Unlock()

	_, _, err := m.ReadFromWebsocket()
	if err != nil {
		return MatchResponse{}, err
	}

	return MatchResponse{
		Message: "Match found",
		MatchDetails: MatchDetails{
			IP:   Localhost,
			Port: fmt.Sprint(DefaultGameserverPort),
		},
	}, nil
}

func (m *MockMatchmaker) CloseMatchmakerConnection(conn *websocket.Conn) error {
	m.Lock()
	defer m.Unlock()
	return m.CloseWebsocketConnection()
}

func (m *MockMatchmaker) DialGameserver(addr *net.UDPAddr) error {
	m.Lock()
	defer m.Unlock()

	m.NumCalledDialGameserver++

	if m.DialGameserverError != nil {
		return m.DialGameserverError
	}

	return nil
}

func (m *MockMatchmaker) CloseGameserverConnection() error {
	m.Lock()
	defer m.Unlock()

	m.NumCalledCloseGameserver++

	if m.CloseGameserverError != nil {
		return m.CloseGameserverError
	}

	return nil
}

func (m *MockMatchmaker) SendMessageToGameserver(message []byte) error {
	m.Lock()
	defer m.Unlock()

	m.NumCalledSendToGameserver++

	if m.SendToGameserverError != nil {
		return m.SendToGameserverError
	}

	return nil
}

func (m *MockMatchmaker) ReadMessageFromGameserver() ([]byte, error) {
	m.Lock()
	defer m.Unlock()

	m.NumCalledReadFromGameserver++

	if m.ReadFromGameserverError != nil {
		return nil, m.ReadFromGameserverError
	}

	return []byte("{\"IPAddress\":\"127.0.0.1\",\"Port\":7777}"), nil
}
