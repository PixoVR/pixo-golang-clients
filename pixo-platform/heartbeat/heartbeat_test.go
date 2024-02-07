package heartbeat_test

import (
	"context"
	graphql_api "github.com/PixoVR/pixo-golang-clients/pixo-platform/graphql-api"
	. "github.com/PixoVR/pixo-golang-clients/pixo-platform/heartbeat"
	"github.com/PixoVR/pixo-golang-clients/pixo-platform/urlfinder"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"os"
)

var _ = Describe("Heartbeat", Ordered, func() {

	var (
		heartbeatClient Client
		config          = urlfinder.ClientConfig{
			Lifecycle: "dev",
		}
		username = os.Getenv("PIXO_USERNAME")
		password = os.Getenv("PIXO_PASSWORD")
	)

	BeforeEach(func() {
		var err error
		heartbeatClient, err = NewClientWithBasicAuth(username, password, config)
		Expect(err).NotTo(HaveOccurred())
		Expect(heartbeatClient).NotTo(BeNil())
		Expect(heartbeatClient.IsAuthenticated()).To(BeTrue())
	})

	It("can create a heartbeatClient and login afterwords", func() {
		anonymousClient := NewClient(config)
		Expect(anonymousClient.IsAuthenticated()).To(BeFalse())
		Expect(anonymousClient.Login(username, password)).NotTo(HaveOccurred())
		Expect(anonymousClient.IsAuthenticated()).To(BeTrue())
	})

	It("should throw an error if the session doesnt exist", func() {
		err := heartbeatClient.SendPulse(9999999999)
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("invalid session"))
	})

	It("should be able to create a session and send a pulse", func() {
		platformClient, err := graphql_api.NewClientWithBasicAuth(username, password, config)
		Expect(err).NotTo(HaveOccurred())
		session, err := platformClient.CreateSession(context.Background(), 271, "127.0.0.1", "test")
		Expect(err).NotTo(HaveOccurred())
		Expect(session).NotTo(BeNil())
		Expect(heartbeatClient.SendPulse(session.ID)).NotTo(HaveOccurred())
	})

})
