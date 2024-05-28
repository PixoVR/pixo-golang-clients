package abstract_client

import (
	"github.com/go-faker/faker/v4"
	"github.com/go-resty/resty/v2"
	"github.com/gorilla/websocket"
	"net/http"
)

type MockAbstractClient struct {
	NumCalledGetIPAddress int
	GetIPAddressError     error

	NumCalledGetURL int

	NumCalledLogin int
	LoginError     error

	NumCalledSetAPIKey int

	NumCalledSetToken int
	Token             string

	NumCalledGetToken int

	NumCalledIsAuthenticated int

	NumCalledActiveUserID int

	NumCalledGet int
	GetError     error

	NumCalledPost int
	PostError     error

	NumCalledPut int
	PutError     error

	NumCalledPatch int
	PatchError     error

	NumCalledDelete int
	DeleteError     error

	NumCalledDialWebsocket int
	DialWebsocketError     error

	NumCalledWriteToWebsocket int
	WriteToWebsocketError     error

	NumCalledReadFromWebsocket int
	ReadFromWebsocketError     error

	NumCalledCloseWebsocket int
	CloseWebsocketError     error

	Response []byte
}

func (m *MockAbstractClient) Path() string {
	return "/api"
}

func (m *MockAbstractClient) GetIPAddress() (string, error) {
	m.NumCalledGetIPAddress++

	if m.GetIPAddressError != nil {
		return "", m.GetIPAddressError
	}

	return "127.0.0.1", nil
}

func (m *MockAbstractClient) GetURL(protocol ...string) string {
	m.NumCalledGetURL++
	return faker.URL()
}

func (m *MockAbstractClient) Login(username, password string) error {
	m.NumCalledLogin++
	if m.LoginError != nil {
		return m.LoginError
	}
	m.SetToken("token")
	return nil
}

func (m *MockAbstractClient) SetAPIKey(key string) {
	m.NumCalledSetAPIKey++
}

func (m *MockAbstractClient) SetToken(token string) {
	m.NumCalledSetToken++
	m.Token = token
}

func (m *MockAbstractClient) GetToken() string {
	m.NumCalledGetToken++
	return m.Token
}

func (m *MockAbstractClient) IsAuthenticated() bool {
	m.NumCalledIsAuthenticated++
	return false
}

func (m *MockAbstractClient) ActiveUserID() int {
	m.NumCalledActiveUserID++
	return 1
}

func (m *MockAbstractClient) Get(path string) (*resty.Response, error) {
	m.NumCalledGet++

	if m.GetError != nil {
		return nil, m.GetError
	}

	return nil, nil
}

func (m *MockAbstractClient) Post(path string, body []byte) (*resty.Response, error) {
	m.NumCalledPost++

	if m.PostError != nil {
		return nil, m.PostError
	}

	return nil, nil
}

func (m *MockAbstractClient) Put(path string, body []byte) (*resty.Response, error) {
	m.NumCalledPut++

	if m.PutError != nil {
		return nil, m.PutError
	}

	return nil, nil
}

func (m *MockAbstractClient) Patch(path string, body []byte) (*resty.Response, error) {
	m.NumCalledPatch++

	if m.PatchError != nil {
		return nil, m.PatchError
	}

	return nil, nil
}

func (m *MockAbstractClient) Delete(path string) (*resty.Response, error) {
	m.NumCalledDelete++

	if m.DeleteError != nil {
		return nil, m.DeleteError
	}

	return nil, nil
}

func (m *MockAbstractClient) DialWebsocket(endpoint string) (*websocket.Conn, *http.Response, error) {
	m.NumCalledDialWebsocket++

	if m.DialWebsocketError != nil {
		return nil, nil, m.DialWebsocketError
	}

	return nil, nil, nil
}

func (m *MockAbstractClient) WriteToWebsocket(message []byte) error {
	m.NumCalledWriteToWebsocket++

	if m.WriteToWebsocketError != nil {
		return m.WriteToWebsocketError
	}

	return nil
}

func (m *MockAbstractClient) ReadFromWebsocket() (int, []byte, error) {
	m.NumCalledReadFromWebsocket++

	if m.ReadFromWebsocketError != nil {
		return 0, nil, m.ReadFromWebsocketError

	}

	return len(m.Response), m.Response, nil
}

func (m *MockAbstractClient) CloseWebsocketConnection() error {
	m.NumCalledCloseWebsocket++

	if m.CloseWebsocketError != nil {
		return m.CloseWebsocketError
	}

	return nil
}
