package matchmaker

import (
	"encoding/json"
	"fmt"
	abstract_client "github.com/PixoVR/pixo-golang-clients/pixo-platform/abstract-client"
	"net"
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
