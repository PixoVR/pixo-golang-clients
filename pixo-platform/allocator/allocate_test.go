package allocator_test

import (
	"context"
	. "github.com/PixoVR/pixo-golang-clients/pixo-platform/allocator"
	"github.com/PixoVR/pixo-golang-clients/pixo-platform/urlfinder"
	config2 "github.com/PixoVR/pixo-golang-server-utilities/pixo-platform/config"
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
			Token:     os.Getenv("TEST_PIXO_SECRET_KEY"),
			Lifecycle: config2.GetEnvOrReturn("TEST_PIXO_LIFECYCLE", "dev"),
		}
		allocatorClient = NewClient(config)
		Expect(allocatorClient.IsAuthenticated()).To(BeTrue())
	})

	It("can check the health of the allocator", func() {
		client := NewClient(config)
		res, err := client.Get(context.Background(), "health")
		Expect(err).NotTo(HaveOccurred())
		Expect(res).NotTo(BeNil())
		Expect(res.StatusCode).To(Equal(http.StatusOK))
	})

	It("should return an error if the server allocation failed", func() {
		allocationReq := AllocationRequest{
			ModuleID:      1,
			OrgID:         1,
			ImageRegistry: "invalid",
		}

		res, err := allocatorClient.AllocateGameserver(allocationReq)

		Expect(err).To(HaveOccurred())
		Expect(res).To(BeNil())
	})

	It("can allocate a multiplayer server", func() {
		req := AllocationRequest{
			ModuleID:           1,
			OrgID:              1,
			ServerVersion:      "1.0.0",
			ImageRegistry:      SimpleGameServerImage,
			AllocateGameServer: true,
		}

		res, err := allocatorClient.AllocateGameserver(req)

		Expect(err).NotTo(HaveOccurred())
		Expect(res).NotTo(BeNil())
		Expect(res.Name).NotTo(BeEmpty())
		Expect(res.IP).NotTo(BeEmpty())
		Expect(res.Port).NotTo(BeEmpty())
		Expect(res.CreatedAt).NotTo(BeEmpty())
	})

})
