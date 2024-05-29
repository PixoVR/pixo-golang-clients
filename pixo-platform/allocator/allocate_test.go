package allocator_test

import (
	. "github.com/PixoVR/pixo-golang-clients/pixo-platform/allocator"
	"github.com/PixoVR/pixo-golang-clients/pixo-platform/urlfinder"
	config2 "github.com/PixoVR/pixo-golang-server-utilities/pixo-platform/config"
	"github.com/PixoVR/pixo-golang-server-utilities/pixo-platform/k8s/agones"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"net/http"
	"os"
)

var _ = Describe("Allocate", Ordered, func() {

	var (
		allocatorClient *Client
		config          urlfinder.ClientConfig
	)

	BeforeEach(func() {
		config = urlfinder.ClientConfig{
			Token:     os.Getenv("SECRET_KEY"),
			Lifecycle: config2.GetEnvOrReturn("PIXO_LIFECYCLE", "stage"),
		}
		allocatorClient = NewClient(config)
		Expect(allocatorClient.IsAuthenticated()).To(BeTrue())
	})

	It("should be able to authenticate with the allocator", func() {
		Expect(allocatorClient.IsAuthenticated()).To(BeTrue())
	})

	It("should be able make a health check against the allocator", func() {
		client := NewClient(config)
		res, err := client.Get("health")
		Expect(err).NotTo(HaveOccurred())
		Expect(res).NotTo(BeNil())
		Expect(res.StatusCode()).To(Equal(http.StatusOK))
	})

	It("should be able to allocate a multiplayer server", func() {
		req := AllocationRequest{
			ModuleID:           1,
			OrgID:              1,
			ServerVersion:      "1.0.0",
			ImageRegistry:      agones.SimpleGameServerImage,
			AllocateGameServer: true,
		}

		res := allocatorClient.AllocateGameserver(req)

		Expect(res.Error).NotTo(HaveOccurred())
		Expect(res.Results.Name).NotTo(BeEmpty())
		Expect(res.Results.IP).NotTo(BeEmpty())
		Expect(res.Results.Port).NotTo(BeEmpty())
		Expect(res.Results.CreatedAt).NotTo(BeEmpty())
	})

	It("should throw an error if the server allocation failed", func() {
		allocationReq := AllocationRequest{
			ModuleID:      1,
			OrgID:         1,
			ImageRegistry: "invalid",
		}
		res := allocatorClient.AllocateGameserver(allocationReq)

		Expect(res).NotTo(BeNil())
		Expect(res.Error).To(HaveOccurred())
	})

})
