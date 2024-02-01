package abstract_client

type AuthClient interface {
	Login(username, password string) error
	SetAPIKey(key string)
	SetToken(key string)
	GetToken() string
	IsAuthenticated() bool
	ActiveUserID() int
}
