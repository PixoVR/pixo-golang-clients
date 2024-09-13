package platform_test

import (
	"context"
	. "github.com/PixoVR/pixo-golang-clients/pixo-platform/platform"
	"github.com/PixoVR/pixo-golang-clients/pixo-platform/urlfinder"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Platform API", func() {

	var (
		platformClient Client
	)

	BeforeEach(func() {
		config := urlfinder.ClientConfig{Lifecycle: lifecycle, APIKey: apiKeyValue}
		platformClient = NewClient(config)
		Expect(platformClient).NotTo(BeNil())
	})

	It("can check if the token is valid", func() {
		user, err := platformClient.CheckAuth(context.Background())
		Expect(err).To(MatchError("unauthorized"))
		Expect(user).NotTo(BeNil())
		Expect(user.ID).To(BeZero())
	})

	It("can login and validate the token", func() {
		Expect(platformClient.Login(username, password)).To(Succeed())
		Expect(platformClient.IsAuthenticated()).To(BeTrue())
		Expect(platformClient.GetToken()).NotTo(BeEmpty())

		user, err := platformClient.CheckAuth(context.Background())
		Expect(err).NotTo(HaveOccurred())
		Expect(user).NotTo(BeNil())
		Expect(user.ID).To(Equal(platformClient.ActiveUserID()))
	})

})
