package platform

import (
	"fmt"
	"github.com/PixoVR/pixo-golang-server-utilities/pixo-platform/middleware/auth"
	"net/http"
)

type transport struct {
	token               string
	key                 string
	underlyingTransport http.RoundTripper
}

func (t *transport) RoundTrip(req *http.Request) (*http.Response, error) {
	if t.key != "" {
		req.Header.Add(auth.APIKeyHeader, t.key)
	} else if t.token != "" {
		req.Header.Add(auth.AuthorizationHeader, fmt.Sprintf("Bearer %s", t.token))
		req.Header.Add(auth.SecretKeyHeader, t.token)
	}

	return t.underlyingTransport.RoundTrip(req)
}
