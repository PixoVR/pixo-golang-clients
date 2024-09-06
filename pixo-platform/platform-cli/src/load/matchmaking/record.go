package matchmaking

import (
	"fmt"
	"github.com/PixoVR/pixo-golang-clients/pixo-platform/matchmaker"
	"github.com/PixoVR/pixo-golang-clients/pixo-platform/platform-cli/src/load/fixture"
	"net"
)

// recordMatchReceived increments the count of successful messages received.
func (t *Tester) recordMatchReceived(response matchmaker.MatchResponse) {
	addr := net.JoinHostPort(response.MatchDetails.IP, response.MatchDetails.Port)
	msg := fmt.Sprintf("gameserver -> %s", addr)

	t.RecordMessage("matchReceived", msg)

	if found := t.gameserverExists(addr); !found {
		t.RecordMessage("gameserverReceived", addr)
	}

	t.Printf(":checkered_flag: %s - %s\n", fixture.SuccessColor.Sprint(response.Message), fixture.StatColor.Sprint(msg))
}

func (t *Tester) gameserverExists(addr string) (found bool) {
	gameservers := t.Messages("gameserverReceived")
	for _, gs := range gameservers {
		if gs == addr {
			found = true
			break
		}
	}
	return found
}

// recordConnectionSuccess increments the count of successful connections.
func (t *Tester) recordConnectionSuccess(id int) {
	t.RecordSuccess(id, "connection", "Connection established")
}

// recordConnectionError adds an error to the list of encountered connectionErrors.
func (t *Tester) recordConnectionError(id int, msg string, err error) {
	t.RecordError(id, "connection", msg, err)
}

// recordMatchingError adds an error to the list of encountered matchingErrors.
func (t *Tester) recordMatchingError(id int, msg string, err error) {
	t.RecordError(id, "match", msg, err)
}
