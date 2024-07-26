package headset_test

import (
	"github.com/PixoVR/pixo-golang-clients/pixo-platform/env"
	"github.com/PixoVR/pixo-golang-clients/pixo-platform/urlfinder"
	"github.com/PixoVR/pixo-golang-server-utilities/pixo-platform/config"
	"os"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestHeadsetClientSuite(t *testing.T) {
	RegisterFailHandler(Fail)
	env.SourceProjectEnv()
	RunSpecs(t, "Headset Client Suite")
}

var (
	lifecycle    string
	username     string
	password     string
	clientConfig urlfinder.ClientConfig
)

var _ = BeforeSuite(func() {
	lifecycle = config.GetEnvOrReturn("TEST_PIXO_LIFECYCLE", "dev")
	username = os.Getenv("TEST_PIXO_USERNAME")
	password = os.Getenv("TEST_PIXO_PASSWORD")
	clientConfig = urlfinder.ClientConfig{Lifecycle: lifecycle}
})
