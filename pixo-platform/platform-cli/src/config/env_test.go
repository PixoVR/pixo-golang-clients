package config_test

import (
	"github.com/PixoVR/pixo-golang-clients/pixo-platform/platform-cli/src/config"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Environment", func() {

	var (
		config = config.Config{}
	)

	It("can get the active environment from a config", func() {
		env := config.ActiveEnv()
		Expect(config.Envs).To(HaveLen(0))
		Expect(env.Lifecycle).To(Equal("prod"))
		Expect(env.Region).To(Equal("na"))
	})

})
