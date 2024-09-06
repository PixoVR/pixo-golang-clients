package heartbeat_test

import (
	"context"
	. "github.com/PixoVR/pixo-golang-clients/pixo-platform/heartbeat"
	"github.com/PixoVR/pixo-golang-clients/pixo-platform/platform"
	"github.com/PixoVR/pixo-golang-clients/pixo-platform/urlfinder"
	config2 "github.com/PixoVR/pixo-golang-server-utilities/pixo-platform/config"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"os"
)

var _ = Describe("Heartbeat Client", Ordered, func() {

	var (
		heartbeatClient Client
		config          = urlfinder.ClientConfig{
			Lifecycle: config2.GetEnvOrReturn("TEST_PIXO_LIFECYCLE", "dev"),
		}
		username = os.Getenv("TEST_PIXO_USERNAME")
		password = os.Getenv("TEST_PIXO_PASSWORD")
		moduleID = 43
	)

	BeforeEach(func() {
		var err error
		heartbeatClient, err = NewClientWithBasicAuth(username, password, config)
		Expect(err).NotTo(HaveOccurred())
		Expect(heartbeatClient).NotTo(BeNil())
		Expect(heartbeatClient.IsAuthenticated()).To(BeTrue())
	})

	It("can create a heartbeat client and login afterwords", func() {
		anonymousClient := NewClient(config)
		Expect(anonymousClient.IsAuthenticated()).To(BeFalse())
		Expect(anonymousClient.Login(username, password)).NotTo(HaveOccurred())
		Expect(anonymousClient.IsAuthenticated()).To(BeTrue())
	})

	It("should return an error if the session doesnt exist", func() {
		err := heartbeatClient.SendPulse(9999999999)
		Expect(err).To(MatchError("invalid session"))
	})

	It("should be able to create a session and send a pulse", func() {
		platformClient, err := platform.NewClientWithBasicAuth(username, password, config)
		Expect(err).NotTo(HaveOccurred())
		session := &platform.Session{
			ModuleID: moduleID,
			DeviceID: "test",
		}
		Expect(platformClient.CreateSession(context.Background(), session)).To(Succeed())
		Expect(err).NotTo(HaveOccurred())
		Expect(session).NotTo(BeNil())

		Expect(heartbeatClient.SendPulse(session.ID)).To(Succeed())
	})

})
