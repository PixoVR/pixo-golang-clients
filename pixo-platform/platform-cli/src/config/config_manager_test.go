package config_test

import (
	"bytes"
	"github.com/PixoVR/pixo-golang-clients/pixo-platform/platform-cli/src/config"
	"github.com/PixoVR/pixo-golang-clients/pixo-platform/platform-cli/src/forms/basic"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/spf13/cobra"
	"os"
)

var _ = Describe("Config Manager", func() {

	var (
		input           *bytes.Buffer
		output          *bytes.Buffer
		formHandler     *basic.Handler
		inMemoryManager *config.InMemoryConfigManager
		configManager   *config.ConfigManager
	)

	BeforeEach(func() {
		input = bytes.NewBufferString("")
		output = bytes.NewBufferString("")
		formHandler = basic.NewFormHandler(input, output)
		inMemoryManager = config.NewInMemoryConfigManager()
		configManager = config.NewConfigManager(inMemoryManager, formHandler)
		Expect(configManager).NotTo(BeNil())
	})

	It("can use the same local config regardless of region", func() {
		Expect(configManager.SetActiveEnv(config.Env{
			Region:    "na",
			Lifecycle: "local",
		})).To(Succeed())

		configManager.SetConfigValue("api-key", "local-api-key")

		Expect(configManager.SetActiveEnv(config.Env{Region: "saudi"})).To(Succeed())
		ExpectConfigValue(configManager, "api-key").To(Equal("local-api-key"))
	})

	It("can set values in the active environment", func() {
		configManager.SetConfigValue("api-key", "new-prod-api-key")
		ExpectConfigValue(configManager, "api-key").To(Equal("new-prod-api-key"))
	})

	It("can prioritize environment variables", func() {
		Expect(os.Setenv("PIXO_REGION", "saudi")).NotTo(HaveOccurred())
		Expect(os.Setenv("PIXO_LIFECYCLE", "dev")).NotTo(HaveOccurred())
		Expect(os.Setenv("PIXO_API_KEY", "saudi-dev-api-key")).NotTo(HaveOccurred())

		Expect(configManager.Region()).To(Equal("saudi"))
		Expect(configManager.Lifecycle()).To(Equal("dev"))
		ExpectConfigValue(configManager, "api-key").To(Equal("saudi-dev-api-key"))

		Expect(os.Unsetenv("PIXO_REGION")).NotTo(HaveOccurred())
		Expect(os.Unsetenv("PIXO_LIFECYCLE")).NotTo(HaveOccurred())
		Expect(os.Unsetenv("PIXO_API_KEY")).NotTo(HaveOccurred())
	})

	It("can get a config value if it exists instead of asking the user", func() {
		configManager.SetConfigValue("username", "na-prod-username")
		username := "new-username"
		input.WriteString(username + "\n")

		val, ok := configManager.GetConfigValueOrAskUser("username", nil)
		Expect(ok).To(BeTrue())
		Expect(val).To(Equal("na-prod-username"))
		Expect(output.String()).To(BeEmpty())
	})

	It("can ask the user if the config doesnt exist", func() {
		configManager.UnsetConfigValue("api-key")
		apiKey := "new-api-key"
		input.WriteString(apiKey + "\n")

		val, ok := configManager.GetConfigValueOrAskUser("api-key", nil)

		Expect(ok).To(BeTrue())
		Expect(val).To(Equal(apiKey))
		Expect(output.String()).To(ContainSubstring("Enter API KEY: "))
	})

	It("can ask the user if the config doesnt exist for an int value", func() {
		configManager.SetIntConfigValue("val", 1)
		val, ok := configManager.GetIntConfigValueOrAskUser("val", nil)
		Expect(ok).To(BeTrue())
		Expect(val).To(Equal(1))
	})

	It("can set an int config value and get it", func() {
		configManager.SetIntConfigValue("port", 8080)
		val, ok := configManager.GetIntConfigValue("port")
		Expect(ok).To(BeTrue())
		Expect(val).To(Equal(8080))
	})

	It("can set a bool config value and get it", func() {
		configManager.SetBoolConfigValue("is-active", true)
		val, ok := configManager.GetBoolConfigValue("is-active")
		Expect(ok).To(BeTrue())
		Expect(val).To(Equal(true))
	})

	It("can prioritize the value from a flag over everything else", func() {
		cmd := &cobra.Command{}
		cmd.Flags().String("api-key", "flag", "api key")
		_ = cmd.Flags().Set("api-key", "flag-api-key")

		val, ok := configManager.GetFlagOrConfigValue("api-key", cmd)
		Expect(ok).To(BeTrue())
		Expect(val).To(Equal("flag-api-key"))
	})

	It("can get an int value from a flag", func() {
		cmd := &cobra.Command{}
		cmd.Flags().Int("port", 8080, "port")
		_ = cmd.Flags().Set("port", "9090")

		val, ok := configManager.GetIntFlagOrConfigValue("port", cmd)
		Expect(ok).To(BeTrue())
		Expect(val).To(Equal(9090))
	})

	It("can get a bool value from a flag", func() {
		cmd := &cobra.Command{}
		cmd.Flags().Bool("is-active", true, "is active")
		_ = cmd.Flags().Set("is-active", "false")

		val, ok := configManager.GetBoolFlagOrConfigValue("is-active", cmd)
		Expect(ok).To(BeTrue())
		Expect(val).To(Equal(false))
	})

})
