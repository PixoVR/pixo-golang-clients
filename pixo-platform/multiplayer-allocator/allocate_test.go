package multiplayer_allocator_test

import (
	. "github.com/PixoVR/pixo-golang-clients/pixo-platform/multiplayer-allocator"
	"github.com/PixoVR/pixo-golang-server-utilities/pixo-platform/k8s/agones"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"net/http"
	"os"
)

var _ = Describe("Allocate", Ordered, func() {

	var (
		allocatorClient *AllocatorClient
	)

	BeforeEach(func() {
		allocatorClient = NewClient(os.Getenv("SECRET_KEY"), "")
		Expect(allocatorClient.IsAuthenticated()).To(BeTrue())
	})

	It("should be able make a health check against the allocator", func() {
		client := NewClient("", "")
		res, err := client.Get("allocator/health")
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
		Expect(res.HTTPResponse.StatusCode()).To(Equal(http.StatusOK))
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

	It("should be able to register a fleet", func() {
		fleetReq := FleetRequest{
			StandbyReplicas: 1,
			ModuleID:        1,
			ImageRegistry:   agones.SimpleGameServerImage,
			ServerVersion:   "1.0.0",
		}

		res := allocatorClient.RegisterFleet(fleetReq)

		Expect(res).NotTo(BeNil())
		Expect(res.Error).NotTo(HaveOccurred())
	})

	It("should be able to deregister a fleet", func() {
		fleetReq := FleetRequest{
			ModuleID:      1,
			ServerVersion: "1.0.0",
		}

		res := allocatorClient.DeregisterFleet(fleetReq)

		Expect(res).NotTo(BeNil())
		Expect(res.Error).NotTo(HaveOccurred())
	})

})
