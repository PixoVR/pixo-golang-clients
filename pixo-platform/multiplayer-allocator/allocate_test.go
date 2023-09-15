package multiplayer_allocator_test

import (
	. "github.com/PixoVR/pixo-golang-clients/pixo-platform/multiplayer-allocator"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"net/http"
	"os"
)

var _ = Describe("Allocate", func() {

	var (
		allocatorClient *AllocatorClient
	)

	BeforeEach(func() {
		allocatorClient = NewClient(os.Getenv("SECRET_KEY"), "")
		Expect(allocatorClient.IsAuthenticated()).To(BeTrue())
	})

	It("should be able make a health check against the allocator", func() {
		client := NewClient("", "")
		res, err := client.Get("health")
		Expect(err).NotTo(HaveOccurred())
		Expect(res.StatusCode()).To(Equal(http.StatusOK))
	})

	It("should be able to allocate a multiplayer server", func() {
		req := AllocateServerRequest{
			ModuleID:           1,
			OrgID:              1,
			ImageRegistry:      "us-docker.pkg.dev/agones-images/examples/simple-game-server:0.14",
			AllocateGameServer: true,
		}

		res, err := allocatorClient.AllocateGameserver(req)

		Expect(err).NotTo(HaveOccurred())
		Expect(res.HTTPResponse.StatusCode()).To(Equal(http.StatusOK))
	})

})
