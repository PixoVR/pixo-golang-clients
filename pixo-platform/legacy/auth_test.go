package legacy_test

import (
	primary_api "github.com/PixoVR/pixo-golang-clients/pixo-platform/legacy"
	"github.com/PixoVR/pixo-golang-clients/pixo-platform/urlfinder"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"os"
)

var _ = Describe("Auth", func() {

	It("should be able to login", func() {
		config := urlfinder.ClientConfig{
			Lifecycle: "dev",
			Region:    "na",
		}
		primaryAPIClient := primary_api.NewClient(config)

		Expect(primaryAPIClient.Login(os.Getenv("TEST_PIXO_USERNAME"), os.Getenv("TEST_PIXO_PASSWORD"))).To(Succeed())

		Expect(primaryAPIClient.IsAuthenticated()).To(BeTrue())
	})

	It("should be able get a new client with basic auth", func() {
		config := urlfinder.ClientConfig{
			Lifecycle: "dev",
			Region:    "na",
		}

		client, err := primary_api.NewClientWithBasicAuth(os.Getenv("TEST_PIXO_USERNAME"), os.Getenv("TEST_PIXO_PASSWORD"), config)

		Expect(err).NotTo(HaveOccurred())
		Expect(client).NotTo(BeNil())
		Expect(client.IsAuthenticated()).To(BeTrue())
		Expect(client.GetToken()).NotTo(BeEmpty())
	})

})
