package abstract

// Response is a struct that represents a basic response from the server.
type Response struct {
	Message string `json:"message"`
	Error   string `json:"error"`
}
