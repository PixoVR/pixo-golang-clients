package primary_api_test

import (
	"github.com/PixoVR/pixo-golang-clients/pixo-platform/urlfinder"
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
		config := urlfinder.ClientConfig{
			Lifecycle: "dev",
			Region:    "na",
			Token:     os.Getenv("SECRET_KEY"),
		}
		primaryClient, err = NewClientWithBasicAuth(os.Getenv("PIXO_USERNAME"), os.Getenv("PIXO_PASSWORD"), config)
		Expect(err).NotTo(HaveOccurred())
		Expect(primaryClient).NotTo(BeNil())
		Expect(primaryClient.IsAuthenticated()).To(BeTrue())

		secretKeyClient = NewClient(config)
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
