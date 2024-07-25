package platform_test

import (
	"context"
	. "github.com/PixoVR/pixo-golang-clients/pixo-platform/platform"
	"github.com/PixoVR/pixo-golang-clients/pixo-platform/urlfinder"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Platform API", func() {

	It("can login and interact with the api", func() {
		config := urlfinder.ClientConfig{Lifecycle: lifecycle, APIKey: pixoAPIKey}
		client := NewClient(config)
		Expect(client).NotTo(BeNil())

		Expect(client.Login(pixoUsername, pixoPassword)).To(Succeed())
		Expect(client.IsAuthenticated()).To(BeTrue())
		Expect(client.GetToken()).NotTo(BeEmpty())

		session, err := client.CreateSession(context.Background(), moduleID, "127.0.0.1", "test")
		Expect(err).NotTo(HaveOccurred())
		Expect(session).NotTo(BeNil())
		Expect(session.ID).NotTo(BeZero())
	})

})
