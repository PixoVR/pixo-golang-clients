package matchmaking

import (
	"errors"
)

// performRequest establishes a single WebSocket connection and requests a match
func (t *Tester) performRequest(id int) {
	conn, _, err := t.config.PlatformFixture.MatchmakingClient.DialMatchmaker()
	if err != nil {
		t.recordConnectionError(id, "failed to connect", err)
		return
	}
	defer func() {
		_ = t.config.PlatformFixture.MatchmakingClient.CloseMatchmakerConnection(conn)
	}()

	t.recordConnectionSuccess(id)

	if err = t.config.PlatformFixture.MatchmakingClient.SendRequest(conn, t.config.Request); err != nil {
		t.recordConnectionError(id, "failed to send request", err)
		return
	}

	matchResponse, err := t.config.PlatformFixture.MatchmakingClient.ReadResponse(conn)
	if err != nil {
		t.recordConnectionError(id, "failed to read message", err)
		return
	}

	if !matchResponse.IsValid() {
		t.recordMatchingError(id, "received invalid match", errors.New(matchResponse.Message))
		return
	}

	t.recordMatchReceived(matchResponse)
}
