package primary_api_test

import (
	"encoding/json"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"net/http"
	"os"

	"github.com/PixoVR/pixo-golang-clients/pixo-platform/primary-api"
)

var _ = Describe("Multiplayer", func() {

	var (
		primaryClient   *primary_api.PrimaryAPIClient
		secretKeyClient *primary_api.PrimaryAPIClient
	)

	BeforeEach(func() {
		primaryClient = primary_api.NewClientWithBasicAuth(os.Getenv("PIXO_USERNAME"), os.Getenv("PIXO_PASSWORD"), "")
		Expect(primaryClient.IsAuthenticated()).To(BeTrue())

		secretKeyClient = primary_api.NewClient(os.Getenv("SECRET_KEY"), "")
		Expect(secretKeyClient.IsAuthenticated()).To(BeTrue())
	})

	It("should be able to update a multiplayer server version using the rest client", func() {
		multiplayerPatch := primary_api.MultiplayerServerVersion{
			Status:        "enabled",
			ImageRegistry: "us-docker.pkg.dev/agones-images/examples/simple-game-server:0.14",
		}
		body, err := json.Marshal(multiplayerPatch)

		res, err := secretKeyClient.Patch("api/external/multiplayer-server-versions/1", body)

		Expect(err).NotTo(HaveOccurred())
		Expect(res.StatusCode()).To(Equal(http.StatusOK))
	})

	It("should be able to update a multiplayer server version using a function", func() {
		image := "us-docker.pkg.dev/agones-images/examples/simple-game-server:0.14"
		res, err := secretKeyClient.UpdateMultiplayerServerVersion(1, image)
		Expect(err).NotTo(HaveOccurred())
		Expect(res.StatusCode()).To(Equal(http.StatusOK))
	})

	It("should be able to create a multiplayer server version using the rest client", func() {
		multiplayerServerVersion := primary_api.MultiplayerServerVersion{
			ModuleID:         17,
			Engine:           "unreal",
			Status:           "enabled",
			ImageRegistry:    "us-docker.pkg.dev/agones-images/examples/simple-game-server:0.14",
			Version:          "1.00.00",
			Filename:         "test.exe",
			MinClientVersion: "1.00.00",
		}
		body, err := json.Marshal(multiplayerServerVersion)

		res, err := primaryClient.Post("api/multiplayer-server-version", body)

		Expect(err).NotTo(HaveOccurred())
		Expect(res).NotTo(BeNil())
		Expect(res.StatusCode()).To(Equal(http.StatusOK))
	})

	It("should be able to deploy a multiplayer server version using a function", func() {
		image := "us-docker.pkg.dev/agones-images/examples/simple-game-server:0.14"
		res, err := primaryClient.DeployMultiplayerServerVersion(17, image, "1.00.00")
		Expect(err).NotTo(HaveOccurred())
		Expect(res.StatusCode()).To(Equal(http.StatusOK))
	})

})
