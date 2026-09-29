package platform

import (
	"context"
	"errors"
	"fmt"
	"net/http"
)

var (
	// ErrNoCredentials is returned when the client has neither an API key nor a token configured.
	ErrNoCredentials = errors.New("no credentials configured: set PIXO_API_KEY or provide a token")

	// ErrUnauthorized is returned when the platform API rejects the credentials the client presents.
	ErrUnauthorized = errors.New("credentials rejected by the platform api")
)

// newResponseError converts a non-200 response from the platform API into an error, wrapping
// ErrUnauthorized when the credentials were rejected so callers can tell that apart from a
// missing record or a malformed request.
func newResponseError(statusCode int, message string) error {
	if statusCode == http.StatusUnauthorized || statusCode == http.StatusForbidden {
		return fmt.Errorf("%w: %s", ErrUnauthorized, message)
	}

	return errors.New(message)
}

// connectionCheckQuery is the cheapest authenticated query available: authentication is enforced
// by middleware before the query is resolved, so it never touches the database.
const connectionCheckQuery = `query { __typename }`

// CheckConnection verifies that the platform API is reachable and that the client's credentials
// are accepted. Call it once after instantiating the client so a missing, revoked or rotated API
// key surfaces immediately instead of as an unrelated-looking failure on every later request.
//
// Callers can distinguish the failure modes with errors.Is(err, ErrNoCredentials) and
// errors.Is(err, ErrUnauthorized); anything else means the API was unreachable.
func (p *clientImpl) CheckConnection(ctx context.Context) error {
	url := p.GetURL()

	if !p.IsAuthenticated() {
		return fmt.Errorf("platform api connection check failed for %s: %w", url, ErrNoCredentials)
	}

	var res struct {
		Typename string `json:"__typename"`
	}

	if err := p.Exec(ctx, connectionCheckQuery, &res, nil); err != nil {
		return fmt.Errorf("platform api connection check failed for %s using the configured %s: %w", url, p.credentialType(), err)
	}

	return nil
}

func (p *clientImpl) credentialType() string {
	if p.GetAPIKey() != "" {
		return "api key"
	}

	return "token"
}
