package abstract_client

import "net/http"

type MockAbstractClient struct {
	NumCalledGetURL          int
	NumCalledLogin           int
	NumCalledSetAPIKey       int
	NumCalledSetToken        int
	NumCalledGetToken        int
	NumCalledIsAuthenticated int
	NumCalledActiveUserID    int

	NumCalledDialWebsocket     int
	NumCalledWriteToWebsocket  int
	NumCalledReadFromWebsocket int
	NumCalledCloseWebsocket    int

	Response []byte
}

func (m *MockAbstractClient) GetURL() string {
	m.NumCalledGetURL++
	return ""
}

func (m *MockAbstractClient) Login(username, password string) error {
	m.NumCalledLogin++
	return nil
}

func (m *MockAbstractClient) SetAPIKey(key string) {
	m.NumCalledSetAPIKey++
}

func (m *MockAbstractClient) SetToken(key string) {
	m.NumCalledSetToken++
}

func (m *MockAbstractClient) GetToken() string {
	m.NumCalledGetToken++
	return ""
}

func (m *MockAbstractClient) IsAuthenticated() bool {
	m.NumCalledIsAuthenticated++
	return false
}

func (m *MockAbstractClient) ActiveUserID() int {
	m.NumCalledActiveUserID++
	return 1
}

func (m *MockAbstractClient) DialWebsocket(endpoint string) (*http.Response, error) {
	m.NumCalledDialWebsocket++
	return nil, nil
}

func (m *MockAbstractClient) WriteToWebsocket(message []byte) error {
	m.NumCalledWriteToWebsocket++
	return nil
}

func (m *MockAbstractClient) ReadFromWebsocket() (int, []byte, error) {
	m.NumCalledReadFromWebsocket++
	return 0, m.Response, nil
}

func (m *MockAbstractClient) CloseWebsocketConnection() error {
	m.NumCalledCloseWebsocket++
	return nil
}
