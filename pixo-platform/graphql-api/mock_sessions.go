package graphql_api

type MockSessionsClient struct {
	CalledGetSession    bool
	CalledCreateSession bool
}

func NewMockSessionsClient() *MockSessionsClient {
	return &MockSessionsClient{}
}

func (m *MockSessionsClient) GetSession(id int) (*Session, error) {
	m.CalledGetSession = true
	return &Session{
		ID:        id,
		UserID:    1,
		ModuleID:  1,
		IPAddress: "127.0.0.1",
		DeviceID:  "1234567890",
	}, nil
}

func (m *MockSessionsClient) CreateSession(moduleID int, ipAddress, deviceId string) (*Session, error) {
	m.CalledCreateSession = true
	return &Session{
		ID:        1,
		UserID:    1,
		ModuleID:  moduleID,
		IPAddress: ipAddress,
		DeviceID:  deviceId,
	}, nil
}
