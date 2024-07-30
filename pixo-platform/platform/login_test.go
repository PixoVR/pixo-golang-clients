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
		config := urlfinder.ClientConfig{Lifecycle: lifecycle, APIKey: apiKey}
		client := NewClient(config)
		Expect(client).NotTo(BeNil())

		Expect(client.Login(username, password)).To(Succeed())
		Expect(client.IsAuthenticated()).To(BeTrue())
		Expect(client.GetToken()).NotTo(BeEmpty())

		platforms, err := client.GetPlatforms(context.Background())
		Expect(err).NotTo(HaveOccurred())
		Expect(platforms).NotTo(BeNil())
		Expect(len(platforms)).To(BeNumerically(">", 0))
	})

})
