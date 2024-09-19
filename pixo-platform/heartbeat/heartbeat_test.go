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
	"time"
)

var _ = Describe("Heartbeat Client", Ordered, func() {

	var (
		ctx = context.Background()

		platformClient  platform.Client
		heartbeatClient Client

		config = urlfinder.ClientConfig{
			Lifecycle: config2.GetEnvOrReturn("TEST_PIXO_LIFECYCLE", "dev"),
		}
		username = os.Getenv("TEST_PIXO_USERNAME")
		password = os.Getenv("TEST_PIXO_PASSWORD")
		moduleID = 43
	)

	BeforeEach(func() {
		var err error
		platformClient, err = platform.NewClientWithBasicAuth(username, password, config)
		Expect(err).NotTo(HaveOccurred())
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
		Expect(heartbeatClient.SendPulse(ctx, -1)).To(MatchError("invalid session"))
	})

	It("can create a session and send a pulse", func() {
		session := &platform.Session{ModuleID: moduleID}
		Expect(platformClient.CreateSession(ctx, session)).To(Succeed())
		Expect(heartbeatClient.SendPulse(ctx, session.ID)).To(Succeed())
	})

	It("sends pulses in a new goroutine", func() {
		errCh, cancel := heartbeatClient.SendPulsesWithCancel(context.Background(), -1, .5)
		Expect(errCh).NotTo(BeNil())
		time.Sleep(1 * time.Second)
		cancel()
		err := <-errCh
		Expect(err).To(MatchError("invalid session"))
		err = <-errCh
		Expect(err).To(MatchError("invalid session"))
	})

})
