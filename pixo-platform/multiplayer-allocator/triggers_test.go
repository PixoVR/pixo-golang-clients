package multiplayer_allocator_test

import (
	. "github.com/PixoVR/pixo-golang-clients/pixo-platform/multiplayer-allocator"
	platform "github.com/PixoVR/pixo-golang-clients/pixo-platform/primary-api"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"net/http"
	"os"
)

var _ = Describe("Triggers", func() {

	var (
		allocatorClient *AllocatorClient
	)

	BeforeEach(func() {
		allocatorClient = NewClient(os.Getenv("SECRET_KEY"), "")
		Expect(allocatorClient.IsAuthenticated()).To(BeTrue())
	})

	It("can register and delete a multiplayer server trigger", func() {
		trigger := platform.MultiplayerServerTrigger{
			ID:         1,
			ModuleID:   1,
			GitOrgName: "PixoVR",
			RepoName:   "multiplayer-gameservers",
			Revision:   "dev",
			Dockerfile: "Server/Dockerfile",
			Context:    ".",
			Config:     "Config/DefaultGame.ini",
		}

		res, err := allocatorClient.RegisterTrigger(trigger)
		Expect(err).NotTo(HaveOccurred())
		Expect(res.StatusCode()).To(Equal(http.StatusCreated))

		res, err = allocatorClient.DeleteTrigger(trigger.ID)
		Expect(err).NotTo(HaveOccurred())
		Expect(res.StatusCode()).To(Equal(http.StatusNoContent))
	})

})
