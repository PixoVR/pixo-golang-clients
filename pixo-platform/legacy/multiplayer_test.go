package legacy_test

import (
	"github.com/PixoVR/pixo-golang-clients/pixo-platform/urlfinder"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"os"

	. "github.com/PixoVR/pixo-golang-clients/pixo-platform/legacy"
)

var _ = Describe("Multiplayer", func() {

	var (
		tokenClient     *Client
		secretKeyClient *Client
	)

	BeforeEach(func() {
		var err error
		config := urlfinder.ClientConfig{
			Lifecycle: os.Getenv("TEST_PIXO_LIFECYCLE"),
			Region:    os.Getenv("TEST_PIXO_REGION"),
		}
		tokenClient, err = NewClientWithBasicAuth(os.Getenv("TEST_PIXO_USERNAME"), os.Getenv("TEST_PIXO_PASSWORD"), config)
		Expect(err).NotTo(HaveOccurred())
		Expect(tokenClient).NotTo(BeNil())
		Expect(tokenClient.IsAuthenticated()).To(BeTrue())

		config.Token = os.Getenv("TEST_PIXO_SECRET_KEY")
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
