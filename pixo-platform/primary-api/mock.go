package primary_api

type MockClient struct {
	NumCalledLogin int
	LoginError     error

	NumCalledCreateWebhook int
	CreateWebhookError     error

	NumCalledGetWebhooks int
	GetWebhooksError     error

	NumCalledDeleteWebhook int
	DeleteWebhookError     error
}

func (m *MockClient) Login(username, password string) error {
	m.NumCalledLogin++
	if m.LoginError != nil {
		return m.LoginError
	}
	return nil
}

func (m *MockClient) CreateWebhook(webhook Webhook) error {
	m.NumCalledCreateWebhook++
	if m.CreateWebhookError != nil {
		return m.CreateWebhookError
	}
	return nil
}

func (m *MockClient) GetWebhooks(orgID int) ([]Webhook, error) {
	m.NumCalledGetWebhooks++
	if m.GetWebhooksError != nil {
		return nil, m.GetWebhooksError
	}
	return []Webhook{
		{
			ID:          1,
			OrgID:       orgID,
			URL:         "https://example.com",
			Description: "test",
		},
	}, nil
}

func (m *MockClient) DeleteWebhook(id int) error {
	m.NumCalledDeleteWebhook++
	if m.DeleteWebhookError != nil {
		return m.DeleteWebhookError
	}
	return nil
}
