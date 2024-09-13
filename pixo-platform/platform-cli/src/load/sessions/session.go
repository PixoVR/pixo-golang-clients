package sessions

import (
	"fmt"
	"github.com/PixoVR/pixo-golang-clients/pixo-platform/headset"
	"github.com/PixoVR/pixo-golang-clients/pixo-platform/platform"
	"github.com/PixoVR/pixo-golang-clients/pixo-platform/platform-cli/src/load/fixture"
)

// performRequest simulates a single session on the platform.
func (t *Tester) performRequest(id int) {
	t.printStart(id)

	if t.config.Legacy {
		t.performLegacySession(id)
	} else {
		t.performSession(id)
	}
}

func (t *Tester) printStart(id int) {
	t.Printf(":checkered_flag: %d: %s\n", id, fixture.StatColor.Sprintf("starting session for module %d...", t.config.Session.ModuleID))
}

func (t *Tester) performSession(id int) {
	session := t.config.Session
	if err := t.config.PlatformFixture.PlatformClient.CreateSession(t.Config.Command.Context(), &session); err != nil {
		t.RecordError(id, "startSession", "unable to start session", err)
		return
	}
	t.RecordSuccess(id, "startSession", fmt.Sprintf("session started for module %s", session.Module.Abbreviation))

	event := platform.Event{
		SessionID: &session.ID,
		Payload:   t.config.Event.Payload,
	}
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

func (t *Tester) performLegacySession(id int) {
	session := t.config.Session
	request := headset.EventRequest{ModuleID: session.ModuleID}
	res, err := t.config.PlatformFixture.HeadsetClient.StartSession(t.Config.Command.Context(), request)
	if err != nil {
		t.RecordError(id, "startSession", "unable to start session", err)
		return
	}
	t.RecordSuccess(id, "startSession", fmt.Sprintf("session started for module %d", session.ModuleID))
	request.SessionID = res.SessionID

	t.preparePayload(&request)
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

func (t *Tester) formatSessionDetails() {
	payload := t.config.Event.Payload
	if payload == nil {
		payload = make(map[string]interface{})
	}

	if t.config.Session.Scenario != "" {
		if payload["object"] == nil {
			payload["object"] = make(map[string]interface{})
		}
		payload["object"].(map[string]interface{})["id"] = fmt.Sprintf("https://pixovr.com/xapi/objects/%d/%s", t.config.Session.ModuleID, t.config.Session.Scenario)
	}

	if t.config.Session.ModuleVersion != "" || t.config.Session.Mode != "" || t.config.Session.Focus != "" || t.config.Session.Specialization != "" {
		if payload["context"] == nil {
			payload["context"] = make(map[string]interface{})
		}

		if t.config.Session.Mode != "" || t.config.Session.Focus != "" || t.config.Session.Specialization != "" {
			payload["context"].(map[string]interface{})["extensions"] = make(map[string]interface{})
		}
	}

	if t.config.Session.ModuleVersion != "" {
		payload["context"].(map[string]interface{})["revision"] = t.config.Session.ModuleVersion
	}

	if t.config.Session.Mode != "" {
		payload["context"].(map[string]interface{})["extensions"].(map[string]interface{})["https://pixovr.com/xapi/extension/sessionMode"] = t.config.Session.Mode
	}

	if t.config.Session.Focus != "" {
		payload["context"].(map[string]interface{})["extensions"].(map[string]interface{})["https://pixovr.com/xapi/extension/sessionFocus"] = t.config.Session.Focus
	}

	if t.config.Session.Specialization != "" {
		payload["context"].(map[string]interface{})["extensions"].(map[string]interface{})["https://pixovr.com/xapi/extension/sessionSpecialization"] = t.config.Session.Specialization
	}

	t.config.SessionDetails = payload
}

func (t *Tester) preparePayload(req *headset.EventRequest) {
	req.Payload = t.config.EventRequest.Payload
	if req.Payload == nil {
		req.Payload = make(map[string]interface{})
	}

	for k, v := range t.config.SessionDetails {
		req.Payload[k] = v
	}
}
