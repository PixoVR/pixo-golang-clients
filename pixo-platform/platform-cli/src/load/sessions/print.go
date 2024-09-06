package sessions

import "github.com/PixoVR/pixo-golang-clients/pixo-platform/platform-cli/src/load/fixture"

// displayStats prints the collected statistics to the console.
func (t *Tester) displayStats() {
	t.Println("Start Session Errors:\t\t" + fixture.ErrorColor.Sprint(len(t.Messages("startSessionError"))))
	t.Println("Create Event Errors:\t\t" + fixture.ErrorColor.Sprint(len(t.Messages("createEventError"))))
	t.Println("Complete Session Errors:\t" + fixture.ErrorColor.Sprint(len(t.Messages("completeSessionError"))))
	t.Println("Unsuccessful Sessions:\t\t" + fixture.ErrorColor.Sprint(t.unsuccessfulSessions()))
	t.Println("Sessions Started:\t\t" + fixture.SuccessColor.Sprint(len(t.Messages("startSessionSuccess"))))
	t.Println("Events Created:\t\t\t" + fixture.SuccessColor.Sprint(len(t.Messages("createEventSuccess"))))
	t.Println("Sessions Completed:\t\t" + fixture.SuccessColor.Sprint(len(t.Messages("completeSessionSuccess"))))
}

func (t *Tester) unsuccessfulSessions() int {
	return len(t.Messages("startSessionError")) +
		len(t.Messages("completeSessionError"))
}
