package sessions

import (
	"fmt"
	"github.com/PixoVR/pixo-golang-clients/pixo-platform/headset"
	"github.com/PixoVR/pixo-golang-clients/pixo-platform/platform"
	"github.com/PixoVR/pixo-golang-clients/pixo-platform/platform-cli/src/load/fixture"
)

// performRequest simulates a single session on the platform.
func (t *Tester) performRequest(id int) {
	session := platform.Session{ModuleID: t.config.Module.ID, Module: t.config.Module}

	t.printStart(id)

	if t.config.Legacy {
		request := headset.EventRequest{ModuleID: session.ModuleID}
		res, err := t.config.PlatformFixture.HeadsetClient.StartSession(t.Config.Command.Context(), request)
		if err != nil {
			t.RecordError(id, "startSession", "unable to start session", err)
			return
		}
		t.RecordSuccess(id, "startSession", fmt.Sprintf("session started for module %d", session.ModuleID))
		request.SessionID = res.SessionID

		request.Payload = t.config.Event.Payload
		if _, err = t.config.PlatformFixture.HeadsetClient.SendEvent(t.Config.Command.Context(), request); err != nil {
			t.RecordError(id, "createEvent", fmt.Sprintf("unable to create event for session %d", request.SessionID), err)
		} else {
			t.RecordSuccess(id, "createEvent", fmt.Sprintf("event created for session %d", request.SessionID))
		}
		request.Payload = nil

		if _, err = t.config.PlatformFixture.HeadsetClient.EndSession(t.Config.Command.Context(), request); err != nil {
			t.RecordError(id, "completeSession", fmt.Sprintf("unable to complete session %d", request.SessionID), err)
			return
		}
		t.RecordSuccess(id, "completeSession", fmt.Sprintf("session completed for module %d", session.ModuleID))

		return
	}

	if err := t.config.PlatformFixture.PlatformClient.CreateSession(t.Config.Command.Context(), &session); err != nil {
		t.RecordError(id, "startSession", "unable to start session", err)
		return
	}
	t.RecordSuccess(id, "startSession", fmt.Sprintf("session started for module %s", session.Module.Abbreviation))

	event := platform.Event{SessionID: &session.ID, Payload: t.config.Event.Payload}
	if err := t.config.PlatformFixture.PlatformClient.CreateEvent(t.Config.Command.Context(), &event); err != nil {
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

func (t *Tester) printStart(id int) {
	t.Printf(":checkered_flag:%d: %s\n", id, fixture.StatColor.Sprintf("starting session for module %d...", t.config.Module.ID))
}
