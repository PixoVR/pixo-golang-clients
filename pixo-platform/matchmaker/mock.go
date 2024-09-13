package matchmaker

import (
	"encoding/json"
	"fmt"
	abstract_client "github.com/PixoVR/pixo-golang-clients/pixo-platform/abstract"
	"github.com/PixoVR/pixo-golang-clients/pixo-platform/allocator"
	"github.com/gorilla/websocket"
	"net"
	"net/http"
)

// MockMatchmaker is a mock implementation of the Matchmaker interface.
type MockMatchmaker struct {
	abstract_client.MockAbstractClient

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

// NewMockMatchmaker creates a new MockMatchmaker instance with a default response.
func NewMockMatchmaker() *MockMatchmaker {
	mockResponse := MatchResponse{
		Message: "Match found",
		MatchDetails: MatchDetails{
			IP:   allocator.Localhost,
			Port: fmt.Sprint(allocator.DefaultGameserverPort),
		},
	}

	mockResponseBytes, _ := json.Marshal(mockResponse)

	return &MockMatchmaker{
		MockAbstractClient: abstract_client.MockAbstractClient{
			Response: mockResponseBytes,
		},
	}
}

// Reset resets the mock's internal state.
func (m *MockMatchmaker) Reset() {
	m.Lock.Lock()
	defer m.Lock.Unlock()

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

// FindMatch is a mock implementation of the Matchmaker interface.
func (m *MockMatchmaker) FindMatch(request MatchRequest) (*net.UDPAddr, error) {
	m.Lock.Lock()
	defer m.Lock.Unlock()

	m.NumCalledFindMatch++

	if m.FindMatchError != nil {
		return nil, m.FindMatchError
	}

	return &net.UDPAddr{
		IP:   net.ParseIP(allocator.Localhost),
		Port: allocator.DefaultGameserverPort,
	}, nil
}

// DialMatchmaker calls the DialWebsocket method on the MockAbstractClient.
func (m *MockMatchmaker) DialMatchmaker() (*websocket.Conn, *http.Response, error) {
	m.Lock.Lock()
	defer m.Lock.Unlock()

	return m.DialWebsocket("")
}

// SendRequest calls the WriteToWebsocket method on the MockAbstractClient with the given request.
func (m *MockMatchmaker) SendRequest(conn *websocket.Conn, req MatchRequest) error {
	m.Lock.Lock()
	defer m.Lock.Unlock()

	reqBytes, _ := json.Marshal(req)
	return m.WriteToWebsocket(reqBytes)
}

// ReadResponse calls the ReadFromWebsocket method on the MockAbstractClient and returns a mock response.
func (m *MockMatchmaker) ReadResponse(conn *websocket.Conn) (MatchResponse, error) {
	m.Lock.Lock()
	defer m.Lock.Unlock()

	_, _, err := m.ReadFromWebsocket()
	if err != nil {
		return MatchResponse{}, err
	}

	return MatchResponse{
		Message: "Match found",
		MatchDetails: MatchDetails{
			IP:   allocator.Localhost,
			Port: fmt.Sprint(allocator.DefaultGameserverPort),
		},
	}, nil
}

// CloseMatchmakerConnection calls the CloseWebsocketConnection method on the MockAbstractClient.
func (m *MockMatchmaker) CloseMatchmakerConnection(conn *websocket.Conn) error {
	m.Lock.Lock()
	defer m.Lock.Unlock()
	return m.CloseWebsocketConnection()
}

// DialGameserver calls the DialGameserver method on the MockAbstractClient.
func (m *MockMatchmaker) DialGameserver(addr *net.UDPAddr) error {
	m.Lock.Lock()
	defer m.Lock.Unlock()

	m.NumCalledDialGameserver++

	if m.DialGameserverError != nil {
		return m.DialGameserverError
	}

	return nil
}

// CloseGameserverConnection returns an error if CloseGameserverError is set.
func (m *MockMatchmaker) CloseGameserverConnection() error {
	m.Lock.Lock()
	defer m.Lock.Unlock()

	m.NumCalledCloseGameserver++

	if m.CloseGameserverError != nil {
		return m.CloseGameserverError
	}

	return nil
}

// SendMessageToGameserver returns an error if SendToGameserverError is set.
func (m *MockMatchmaker) SendMessageToGameserver(message []byte) error {
	m.Lock.Lock()
	defer m.Lock.Unlock()

	m.NumCalledSendToGameserver++

	if m.SendToGameserverError != nil {
		return m.SendToGameserverError
	}

	return nil
}

// ReadMessageFromGameserver returns an error if ReadFromGameserverError is set, if not is returns a mock response string.
func (m *MockMatchmaker) ReadMessageFromGameserver() ([]byte, error) {
	m.Lock.Lock()
	defer m.Lock.Unlock()

	m.NumCalledReadFromGameserver++

	if m.ReadFromGameserverError != nil {
		return nil, m.ReadFromGameserverError
	}

	return []byte("{\"IPAddress\":\"127.0.0.1\",\"Port\":7777}"), nil
}
