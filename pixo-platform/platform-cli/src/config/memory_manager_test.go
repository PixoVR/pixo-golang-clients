package config_test

import (
	"github.com/PixoVR/pixo-golang-clients/pixo-platform/platform-cli/src/config"
	"github.com/PixoVR/pixo-golang-clients/pixo-platform/platform-cli/src/printer"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("In Memory Manager", func() {

	var (
		inMemoryConfigManager *config.InMemoryConfigManager
		configManager         *config.ConfigManager
	)

	BeforeEach(func() {
		inMemoryConfigManager = config.NewInMemoryConfigManager()
		Expect(inMemoryConfigManager).NotTo(BeNil())
		emojiPrinter := printer.NewEmojiPrinter(nil)
		configManager = config.NewConfigManager(inMemoryConfigManager, emojiPrinter)
		Expect(configManager).NotTo(BeNil())
	})

	It("can get the active environment", func() {
		env := configManager.ActiveEnv()
		Expect(env.Lifecycle).To(Equal("prod"))
		Expect(env.Region).To(Equal("na"))
	})

	It("can set the active environment", func() {
		Expect(configManager.SetActiveEnv(config.Env{
			Region:    "na",
			Lifecycle: "stage",
		})).To(Succeed())

		env := configManager.ActiveEnv()

		Expect(env.Region).To(Equal("na"))
		Expect(env.Lifecycle).To(Equal("stage"))
		Expect(env.EnvMap).NotTo(BeNil())
	})

	It("can set a config value", func() {
		apiKey := "some-api-key"
		configManager.SetConfigValue("api-key", apiKey)
		ExpectConfigValue(configManager, "api-key").To(Equal(apiKey))
	})

	It("can set values in the active environment", func() {
		configManager.SetConfigValue("api-key", "new-prod-api-key")
		ExpectConfigValue(configManager, "api-key").To(Equal("new-prod-api-key"))
	})

})
