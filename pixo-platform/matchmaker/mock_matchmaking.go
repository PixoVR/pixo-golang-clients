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
	NumCalledFindMatch       int
	NumCalledDialGameserver  int
	NumCalledSendAndReceive  int
	NumCalledCloseGameserver int
	response                 []byte
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

	return &net.UDPAddr{
		IP:   net.ParseIP(Localhost),
		Port: DefaultGameserverPort,
	}, nil
}

func (m *MockMatchmaker) DialMatchmaker() (*websocket.Conn, *http.Response, error) {
	m.NumCalledDialWebsocket++
	return nil, nil, nil
}

func (m *MockMatchmaker) SendRequest(conn *websocket.Conn, req MatchRequest) error {
	m.NumCalledWriteToWebsocket++
	return nil
}

func (m *MockMatchmaker) ReadResponse(conn *websocket.Conn) (MatchResponse, error) {
	m.NumCalledReadFromWebsocket++
	response := MatchResponse{}
	err := json.Unmarshal(m.response, &response)
	if err != nil {
		return MatchResponse{}, err
	}

	return response, nil
}

func (m *MockMatchmaker) CloseMatchmakerConnection(conn *websocket.Conn) error {
	m.NumCalledCloseWebsocket++
	return nil
}

func (m *MockMatchmaker) DialGameserver(addr *net.UDPAddr) error {
	m.NumCalledDialGameserver++
	return nil
}

func (m *MockMatchmaker) CloseGameserverConnection() error {
	m.NumCalledCloseGameserver++
	return nil
}

func (m *MockMatchmaker) SendAndReceiveMessage(message []byte) ([]byte, error) {
	m.NumCalledSendAndReceive++
	return []byte("{\"IPAddress\":\"127.0.0.1\",\"Port\":7777}"), nil
}
