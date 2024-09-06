package matchmaking

import (
	"github.com/PixoVR/pixo-golang-clients/pixo-platform/platform-cli/src/load/fixture"
)

// displayStats prints the collected statistics to the console.
func (t *Tester) displayStats() {
	t.Println(fixture.ErrorColor.Sprintf("Connection Errors:         %d", len(t.Messages("connectError"))))
	t.Println(fixture.SuccessColor.Sprintf("Successful Connections:    %d", len(t.Messages("matchReceived"))))
	t.Println(fixture.ErrorColor.Sprintf("Matching Errors:           %d", len(t.Messages("matchError"))))
	t.Println(fixture.SuccessColor.Sprintf("Matches Received:          %d", len(t.Messages("matchReceived"))))
	t.Println(fixture.SuccessColor.Sprintf("Gameservers Received:      %d", len(t.Messages("gameserverReceived"))))
}
