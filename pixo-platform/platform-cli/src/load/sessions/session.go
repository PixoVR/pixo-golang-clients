package sessions

import (
	"fmt"
	"github.com/PixoVR/pixo-golang-clients/pixo-platform/platform"
)

// performRequest simulates a single session on the platform.
func (t *Tester) performRequest(id int) {
	session := platform.Session{ModuleID: t.config.Module.ID, Module: t.config.Module}
	if err := t.config.PlatformFixture.PlatformClient.CreateSession(t.Config.Command.Context(), &session); err != nil {
		t.RecordError(id, "startSession", "unable to start session", err)
		return
	}
	t.RecordSuccess(id, "startSession", fmt.Sprintf("session started for module %s", session.Module.Abbreviation))

	if err := t.config.PlatformFixture.PlatformClient.CreateEvent(t.Config.Command.Context(), &platform.Event{SessionID: &session.ID}); err != nil {
		t.RecordError(id, "createEvent", "unable to create event", err)
	} else {
		t.RecordSuccess(id, "createEvent", fmt.Sprintf("event created for session %d", session.ID))
	}

	if _, err := t.config.PlatformFixture.PlatformClient.UpdateSession(t.Config.Command.Context(), session); err != nil {
		t.RecordError(id, "completeSession", "unable to complete session", err)
		return
	}
	t.RecordSuccess(id, "completeSession", fmt.Sprintf("session completed for module %s", session.Module.Abbreviation))
}
