package abstract_client

type AbstractClient interface {
	Login(username, password string) error
	SetAPIKey(key string)
	SetToken(key string)
	GetToken() string
	GetURL() string
	IsAuthenticated() bool
	ActiveUserID() int
}
