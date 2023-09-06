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
		primaryClient *primary_api.PrimaryAPIClient
	)

	BeforeEach(func() {
		primaryClient = primary_api.NewClient(os.Getenv("SECRET_KEY"), "")
		Expect(primaryClient.IsAuthenticated()).To(BeTrue())
	})

	It("should be able to update a multiplayer server version using the rest client", func() {
		multiplayerPatch := primary_api.MultiplayerServerVersion{
			Status:        "enabled",
			ImageRegistry: "us-docker.pkg.dev/agones-images/examples/simple-game-server:0.14",
		}
		body, err := json.Marshal(multiplayerPatch)
		res, err := primaryClient.Patch("api/external/multiplayer-server-versions/1", body)
		Expect(err).NotTo(HaveOccurred())
		Expect(res.StatusCode()).To(Equal(http.StatusOK))
	})

	It("should be able to update a multiplayer server version using a function", func() {
		imageRegistry := "us-docker.pkg.dev/agones-images/examples/simple-game-server:0.14"
		status := "enabled"
		res, err := primaryClient.UpdateMultiplayerServerVersion(imageRegistry, status)
		Expect(err).NotTo(HaveOccurred())
		Expect(res.StatusCode()).To(Equal(http.StatusOK))
	})

})
