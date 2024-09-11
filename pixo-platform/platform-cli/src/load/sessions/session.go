package sessions

import (
	"fmt"
	"github.com/PixoVR/pixo-golang-clients/pixo-platform/headset"
	"github.com/PixoVR/pixo-golang-clients/pixo-platform/platform"
)

// performRequest simulates a single session on the platform.
func (t *Tester) performRequest(id int) {
	session := platform.Session{ModuleID: t.config.Module.ID, Module: t.config.Module}

	if t.config.Legacy {
		request := headset.EventRequest{ModuleID: session.ModuleID}
		res, err := t.config.PlatformFixture.HeadsetClient.StartSession(t.Config.Command.Context(), request)
		if err != nil {
			t.RecordError(id, "startSession", "unable to start session", err)
			return
		}
		t.RecordSuccess(id, "startSession", fmt.Sprintf("session started for module %d", t.config.Module.ID))
		request.SessionID = res.SessionID

		if _, err = t.config.PlatformFixture.HeadsetClient.SendEvent(t.Config.Command.Context(), request); err != nil {
			t.RecordError(id, "createEvent", fmt.Sprintf("unable to create event for session %d", request.SessionID), err)
		} else {
			t.RecordSuccess(id, "createEvent", fmt.Sprintf("event created for session %d", request.SessionID))
		}

		if _, err = t.config.PlatformFixture.HeadsetClient.EndSession(t.Config.Command.Context(), request); err != nil {
			t.RecordError(id, "completeSession", fmt.Sprintf("unable to complete session %d", request.SessionID), err)
			return
		}
		t.RecordSuccess(id, "completeSession", fmt.Sprintf("session completed for module %d", t.config.Module.ID))

		return
	}

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

	session.Completed = true
	if _, err := t.config.PlatformFixture.PlatformClient.UpdateSession(t.Config.Command.Context(), session); err != nil {
		t.RecordError(id, "completeSession", "unable to complete session", err)
		return
	}
	t.RecordSuccess(id, "completeSession", fmt.Sprintf("session completed for module %s", session.Module.Abbreviation))
}
