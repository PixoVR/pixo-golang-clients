package platform_test

import (
	"context"
	"github.com/PixoVR/pixo-golang-clients/pixo-platform/platform"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Sessions and Events", func() {

	var (
		ctx       = context.Background()
		session   *platform.Session
		ipAddress = "127.0.0.1"
		deviceID  = "test"
	)

	BeforeEach(func() {
		session = &platform.Session{
			ModuleID:  moduleID,
			IPAddress: ipAddress,
			DeviceID:  deviceID,
		}

		Expect(tokenClient.CreateSession(ctx, session)).To(Succeed())

		Expect(session).NotTo(BeNil())
		Expect(session.ID).NotTo(BeZero())
		Expect(session.UserID).NotTo(BeZero())
		Expect(session.ModuleID).To(Equal(moduleID))
		Expect(session.Module).NotTo(BeNil())
		Expect(session.Module.ID).To(Equal(moduleID))
		Expect(session.User).NotTo(BeNil())
		Expect(session.User.OrgID).NotTo(BeZero())
	})

	It("can get a session", func() {
		retrievedSession, err := tokenClient.GetSession(ctx, session.ID)

		Expect(err).NotTo(HaveOccurred())
		Expect(retrievedSession).NotTo(BeNil())
		Expect(retrievedSession.ID).To(Equal(session.ID))
		Expect(retrievedSession.UserID).NotTo(BeZero())
	})

	It("can update a session", func() {
		input := platform.Session{
			ID:        session.ID,
			Status:    "TERMINATED",
			Completed: true,
			RawScore:  0.5,
			MaxScore:  1.0,
		}

		updatedSession, err := tokenClient.UpdateSession(ctx, input)

		Expect(err).NotTo(HaveOccurred())
		Expect(updatedSession).NotTo(BeNil())
		Expect(updatedSession.ID).To(Equal(session.ID))
		Expect(updatedSession.UserID).To(Equal(session.UserID))
		Expect(updatedSession.UserID).To(Equal(session.UserID))
		Expect(updatedSession.User.OrgID).To(Equal(session.User.OrgID))
		Expect(updatedSession.ModuleID).To(Equal(session.ModuleID))
		Expect(updatedSession.RawScore).To(Equal(input.RawScore))
		Expect(updatedSession.MaxScore).To(Equal(input.MaxScore))
		Expect(updatedSession.ScaledScore).To(BeNumerically("~", input.RawScore/input.MaxScore, 0.01))
		Expect(updatedSession.CompletedAt).NotTo(BeNil())
		Expect(updatedSession.Duration).NotTo(BeNil())
	})

	It("can return an error if the json is invalid", func() {
		event := &platform.Event{
			SessionID: session.ID,
			Type:      "PIXOVR_SESSION_JOIN",
			Payload:   `{"missing": "end bracket"`,
		}

		err := tokenClient.CreateEvent(ctx, event)

		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("invalid json"))
	})

	It("can create an event without a payload", func() {
		event := &platform.Event{
			SessionID: session.ID,
			Type:      "PIXOVR_SESSION_JOIN",
		}

		Expect(tokenClient.CreateEvent(ctx, event)).To(Succeed())

		Expect(event).NotTo(BeNil())
		Expect(event.ID).NotTo(BeZero())
		Expect(event.SessionID).To(Equal(session.ID))
	})

	It("can create an event with a payload", func() {
		event := &platform.Event{
			SessionID: session.ID,
			Type:      "PIXOVR_SESSION_JOIN",
			Payload:   `{"score": 1}`,
		}

		Expect(tokenClient.CreateEvent(ctx, event)).To(Succeed())

		Expect(event).NotTo(BeNil())
		Expect(event.ID).NotTo(BeZero())
	})

})
