package primary_api_test

import (
	primary_api "github.com/PixoVR/pixo-golang-clients/pixo-platform/primary-api"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"os"
)

var _ = Describe("Auth", func() {

	It("should be able to login", func() {
		primaryAPIClient := primary_api.NewClient("", "dev", "na")
		err := primaryAPIClient.Login(os.Getenv("PIXO_USERNAME"), os.Getenv("PIXO_PASSWORD"))
		Expect(err).NotTo(HaveOccurred())
		Expect(primaryAPIClient.IsAuthenticated()).To(BeTrue())
	})

	It("should be able get a new client with basic auth", func() {
		client, err := primary_api.NewClientWithBasicAuth(os.Getenv("PIXO_USERNAME"), os.Getenv("PIXO_PASSWORD"), "dev", "")
		Expect(err).NotTo(HaveOccurred())
		Expect(client).NotTo(BeNil())
		Expect(client.IsAuthenticated()).To(BeTrue())
		Expect(client.GetToken()).NotTo(BeEmpty())
	})

})
