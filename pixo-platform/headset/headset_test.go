package headset_test

import (
	"context"
	. "github.com/PixoVR/pixo-golang-clients/pixo-platform/headset"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Headset Client", func() {

	var (
		headsetClient Client
		ctx           = context.Background()
	)

	BeforeEach(func() {
		var err error
		headsetClient, err = NewClientWithBasicAuth(username, password, clientConfig)
		Expect(err).NotTo(HaveOccurred())
		Expect(headsetClient).NotTo(BeNil())
		Expect(headsetClient.IsAuthenticated()).To(BeTrue())
	})

	It("can login", func() {
		anonymousClient := NewClient(clientConfig)
		Expect(anonymousClient.IsAuthenticated()).To(BeFalse())
		Expect(anonymousClient.Login(username, password)).NotTo(HaveOccurred())
		Expect(anonymousClient.IsAuthenticated()).To(BeTrue())
	})

	It("can perform a session", func() {
		input := EventRequest{
			ModuleID: moduleID,
		}

		response, err := headsetClient.StartSession(ctx, input)

		Expect(err).NotTo(HaveOccurred())
		Expect(response).NotTo(BeNil())
		Expect(response.SessionID).NotTo(BeZero())

		input.SessionID = response.SessionID
		input.Payload = map[string]interface{}{
			"some": "important data",
		}

		response, err = headsetClient.SendEvent(ctx, input)

		Expect(err).NotTo(HaveOccurred())
		Expect(response).NotTo(BeNil())
		Expect(response.SessionID).NotTo(BeZero())

		input.Payload = map[string]interface{}{
			"score":        190,
			"scoreMax":     200,
			"lessonStatus": "passed",
		}

		response, err = headsetClient.EndSession(ctx, input)

		Expect(err).NotTo(HaveOccurred())
		Expect(response).NotTo(BeNil())
		Expect(response.SessionID).NotTo(BeZero())
	})

})
