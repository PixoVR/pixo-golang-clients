package abstract_client

type Response struct {
	Message string `json:"message"`
	Error   string `json:"error"`
}
