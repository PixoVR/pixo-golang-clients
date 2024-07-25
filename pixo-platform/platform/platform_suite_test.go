package platform_test

import (
	"github.com/PixoVR/pixo-golang-clients/pixo-platform/env"
	"github.com/PixoVR/pixo-golang-clients/pixo-platform/urlfinder"
	config2 "github.com/PixoVR/pixo-golang-server-utilities/pixo-platform/config"
	"os"
	"testing"

	. "github.com/PixoVR/pixo-golang-clients/pixo-platform/platform"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestPlatformClientSuite(t *testing.T) {
	RegisterFailHandler(Fail)
	env.SourceProjectEnv()
	RunSpecs(t, "Platform Client Suite")
}

var (
	apiKeyClient *PlatformClient
	tokenClient  *PlatformClient
	lifecycle    string
	pixoUsername string
	pixoPassword string
	pixoAPIKey   string
	moduleID     = 43
	orgID        = 20
)

var _ = BeforeSuite(func() {
	lifecycle = config2.GetEnvOrReturn("TEST_PIXO_LIFECYCLE", "dev")
	pixoUsername = os.Getenv("TEST_PIXO_USERNAME")
	pixoPassword = os.Getenv("TEST_PIXO_PASSWORD")
	pixoAPIKey = os.Getenv("TEST_PIXO_API_KEY")

	config := urlfinder.ClientConfig{Lifecycle: lifecycle, APIKey: pixoAPIKey}
	apiKeyClient = NewClient(config)
	Expect(apiKeyClient).NotTo(BeNil())
	Expect(apiKeyClient.IsAuthenticated()).To(BeTrue())

	var err error
	tokenClient, err = NewClientWithBasicAuth(pixoUsername, pixoPassword, config)
	Expect(err).NotTo(HaveOccurred())
	Expect(tokenClient).NotTo(BeNil())
	Expect(tokenClient.IsAuthenticated()).To(BeTrue())
	Expect(tokenClient.GetToken()).NotTo(BeEmpty())
})
