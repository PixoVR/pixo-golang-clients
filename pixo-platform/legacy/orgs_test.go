package legacy_test

import (
	"github.com/PixoVR/pixo-golang-clients/pixo-platform/legacy"
	"github.com/PixoVR/pixo-golang-clients/pixo-platform/urlfinder"
	config2 "github.com/PixoVR/pixo-golang-server-utilities/pixo-platform/config"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"os"
)

var _ = Describe("Orgs", Ordered, func() {

	var (
		apiClient *legacy.Client
	)

	BeforeEach(func() {
		config := urlfinder.ClientConfig{
			Lifecycle: config2.GetEnvOrReturn("TEST_PIXO_LIFECYCLE", "dev"),
			Region:    config2.GetEnvOrReturn("TEST_PIXO_REGION", "na"),
		}
		apiClient = legacy.NewClient(config)
		Expect(apiClient.Login(os.Getenv("TEST_PIXO_USERNAME"), os.Getenv("TEST_PIXO_PASSWORD"))).To(Succeed())
	})

	It("can get orgs", func() {
		orgs, err := apiClient.GetOrgs()
		Expect(err).NotTo(HaveOccurred())
		Expect(orgs).NotTo(BeNil())
		Expect(len(orgs)).To(BeNumerically(">", 0))
	})

})
