package primary_api_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"net/http"
	"os"

	. "github.com/PixoVR/pixo-golang-clients/pixo-platform/primary-api"
)

var _ = Describe("Multiplayer", func() {

	var (
		primaryClient   *PrimaryAPIClient
		secretKeyClient *PrimaryAPIClient
	)

	BeforeEach(func() {
		primaryClient = NewClientWithBasicAuth(os.Getenv("PIXO_USERNAME"), os.Getenv("PIXO_PASSWORD"), "")
		Expect(primaryClient).NotTo(BeNil())
		Expect(primaryClient.IsAuthenticated()).To(BeTrue())

		secretKeyClient = NewClient(os.Getenv("SECRET_KEY"), "")
		Expect(secretKeyClient).NotTo(BeNil())
		Expect(secretKeyClient.IsAuthenticated()).To(BeTrue())
	})

	It("should be able to update a multiplayer server version", func() {
		image := "us-docker.pkg.dev/agones-images/examples/simple-game-server:0.14"
		res, err := secretKeyClient.UpdateMultiplayerServerVersion(1, image)
		Expect(err).NotTo(HaveOccurred())
		Expect(res.StatusCode()).To(Equal(http.StatusOK))

	})

	It("should be able to get the multiplayer configurations available", func() {
		profiles, err := secretKeyClient.GetMatchmakingProfiles()
		Expect(err).NotTo(HaveOccurred())
		Expect(profiles).NotTo(BeNil())
		Expect(len(profiles)).To(BeNumerically(">", 0))
	})

	It("should be able to deploy a multiplayer server version", func() {
		image := "us-docker.pkg.dev/agones-images/examples/simple-game-server:0.14"
		res, err := primaryClient.DeployMultiplayerServerVersion(125, image, "1.00.00")
		Expect(err).NotTo(HaveOccurred())
		Expect(res.StatusCode()).To(Equal(http.StatusOK))
	})

})
