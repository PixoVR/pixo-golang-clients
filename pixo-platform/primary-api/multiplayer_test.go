package primary_api_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"os"

	. "github.com/PixoVR/pixo-golang-clients/pixo-platform/primary-api"
)

var _ = Describe("Multiplayer", func() {

	var (
		primaryClient   *PrimaryAPIClient
		secretKeyClient *PrimaryAPIClient
	)

	BeforeEach(func() {
		var err error
		primaryClient, err = NewClientWithBasicAuth(os.Getenv("PIXO_USERNAME"), os.Getenv("PIXO_PASSWORD"), "dev", "na")
		Expect(err).NotTo(HaveOccurred())
		Expect(primaryClient).NotTo(BeNil())
		Expect(primaryClient.IsAuthenticated()).To(BeTrue())

		secretKeyClient = NewClient(os.Getenv("SECRET_KEY"), "dev", "na")
		Expect(secretKeyClient).NotTo(BeNil())
		Expect(secretKeyClient.IsAuthenticated()).To(BeTrue())
	})

	It("should be able to get the matchmaking profiles available", func() {
		profiles, err := secretKeyClient.GetMatchmakingProfiles()
		Expect(err).NotTo(HaveOccurred())
		Expect(profiles).NotTo(BeNil())
		Expect(len(profiles)).To(BeNumerically(">", 0))
	})

})
