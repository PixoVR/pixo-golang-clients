package abstract

import (
	"context"
	"github.com/PixoVR/pixo-golang-server-utilities/pixo-platform/config"
	"github.com/go-faker/faker/v4"
	"github.com/golang-jwt/jwt/v4"
	"github.com/gorilla/websocket"
	"github.com/rs/zerolog/log"
	"io"
	"net/http"
	"sync"
)

var _ AbstractClient = (*MockAbstractClient)(nil)

// MockAbstractClient is a mock implementation of the AbstractClient interface
type MockAbstractClient struct {
	Lock sync.Mutex

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

	NumCalledWriteToWebsocketError int
	WriteToWebsocketError          error

	NumCalledReadFromWebsocket int
	ReadFromWebsocketError     error

	NumCalledCloseWebsocket int
	CloseWebsocketError     error

	Response []byte
}

// Path returns "/api"
func (m *MockAbstractClient) Path() string {
	return "/api"
}

// GetIPAddress returns an error if provided or localhost
func (m *MockAbstractClient) GetIPAddress() (string, error) {
	m.NumCalledGetIPAddress++

	if m.GetIPAddressError != nil {
		return "", m.GetIPAddressError
	}

	return "127.0.0.1", nil
}

// GetURL returns a random URL
func (m *MockAbstractClient) GetURL(protocol ...string) string {
	m.NumCalledGetURL++
	return faker.URL()
}

// NewRequest returns a new request
func (m *MockAbstractClient) NewRequest(method, path string, body io.Reader) (*http.Request, error) {
	return http.NewRequest(method, path, nil)
}

// Login sets the token to a fake JWT
func (m *MockAbstractClient) Login(username, password string) error {
	m.NumCalledLogin++
	if m.LoginError != nil {
		return m.LoginError
	}

	claims := jwt.MapClaims{
		"authorized": true,
		"userId":     1,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(config.GetEnvOrReturn("SECRET_KEY", "fake-key")))
	if err != nil {
		log.Panic().Err(err).Msg("error signing token")
	}

	m.SetToken(signedToken)
	return nil
}

// SetAPIKey increments the number of times it was called
func (m *MockAbstractClient) SetAPIKey(key string) {
	m.NumCalledSetAPIKey++
}

// SetToken sets the token to the provided value
func (m *MockAbstractClient) SetToken(token string) {
	m.NumCalledSetToken++
	m.Token = token
}

// GetToken returns the token
func (m *MockAbstractClient) GetToken() string {
	m.NumCalledGetToken++
	return m.Token
}

// IsAuthenticated returns false
func (m *MockAbstractClient) IsAuthenticated() bool {
	m.NumCalledIsAuthenticated++
	return false
}

// ActiveUserID returns 1
func (m *MockAbstractClient) ActiveUserID() int {
	m.NumCalledActiveUserID++
	return 1
}

// Get increments the number of times it was called and returns an error if provided
func (m *MockAbstractClient) Get(ctx context.Context, path string) (*http.Response, error) {
	m.NumCalledGet++

	if m.GetError != nil {
		return nil, m.GetError
	}

	return nil, nil
}

// Post increments the number of times it was called and returns an error if provided
func (m *MockAbstractClient) Post(ctx context.Context, path string, body []byte) (*http.Response, error) {
	m.NumCalledPost++

	if m.PostError != nil {
		return nil, m.PostError
	}

	return nil, nil
}

// Put increments the number of times it was called and returns an error if provided
func (m *MockAbstractClient) Put(ctx context.Context, path string, body []byte) (*http.Response, error) {
	m.NumCalledPut++

	if m.PutError != nil {
		return nil, m.PutError
	}

	return nil, nil
}

// Patch increments the number of times it was called and returns an error if provided
func (m *MockAbstractClient) Patch(ctx context.Context, path string, body []byte) (*http.Response, error) {
	m.NumCalledPatch++

	if m.PatchError != nil {
		return nil, m.PatchError
	}

	return nil, nil
}

// Delete increments the number of times it was called and returns an error if provided
func (m *MockAbstractClient) Delete(ctx context.Context, path string) (*http.Response, error) {
	m.NumCalledDelete++

	if m.DeleteError != nil {
		return nil, m.DeleteError
	}

	return nil, nil
}

// DialWebsocket increments the number of times it was called and returns an error if provided
func (m *MockAbstractClient) DialWebsocket(endpoint string) (*websocket.Conn, *http.Response, error) {
	m.NumCalledDialWebsocket++

	if m.DialWebsocketError != nil {
		return nil, nil, m.DialWebsocketError
	}

	return nil, nil, nil
}

// WriteToWebsocket increments the number of times it was called and returns an error if provided
func (m *MockAbstractClient) WriteToWebsocket(message []byte) error {
	m.NumCalledWriteToWebsocketError++

	if m.WriteToWebsocketError != nil {
		return m.WriteToWebsocketError
	}

	return nil
}

// ReadFromWebsocket increments the number of times it was called and returns an error if provided
func (m *MockAbstractClient) ReadFromWebsocket() (int, []byte, error) {
	m.NumCalledReadFromWebsocket++

	if m.ReadFromWebsocketError != nil {
		return 0, nil, m.ReadFromWebsocketError

	}

	return len(m.Response), m.Response, nil
}

// CloseWebsocketConnection increments the number of times it was called and returns an error if provided
func (m *MockAbstractClient) CloseWebsocketConnection() error {
	m.NumCalledCloseWebsocket++

	if m.CloseWebsocketError != nil {
		return m.CloseWebsocketError
	}

	return nil
}
