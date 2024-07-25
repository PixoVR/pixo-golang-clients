package config_test

import (
	"bytes"
	"fmt"
	"github.com/PixoVR/pixo-golang-clients/pixo-platform/platform-cli/src/config"
	"github.com/PixoVR/pixo-golang-clients/pixo-platform/platform-cli/src/forms/basic"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"math/rand"
	"os"
)

var _ = Describe("File Manager", func() {

	Context("when the config file does not exist", func() {

		var (
			filename      string
			fileManager   *config.FileManager
			configManager *config.ConfigManager
		)

		BeforeEach(func() {
			filename = fmt.Sprintf("./test-config-%d.yaml", rand.Intn(100000))
			fileManager = config.NewFileConfigManager(filename)
			Expect(fileManager).NotTo(BeNil())
			configManager = config.NewConfigManager(fileManager)
			Expect(configManager).NotTo(BeNil())
		})

		AfterEach(func() {
			viper.Reset()
			err := os.Remove(filename)
			Expect(err).NotTo(HaveOccurred())
		})

		It("can initialize the config manager when setting the active environment", func() {
			Expect(configManager.SetActiveEnv(config.Env{
				Region:    "na",
				Lifecycle: "stage",
			})).To(Succeed())

			env := configManager.ActiveEnv()

			Expect(env.Region).To(Equal("na"))
			Expect(env.Lifecycle).To(Equal("stage"))
		})

		It("can initialize the config manager when setting a value", func() {
			apiKey := "some-api-key"
			configManager.SetConfigValue("api-key", apiKey)
			ExpectConfigValue(configManager, "api-key").To(Equal(apiKey))
		})

		It("can return an error if the region does not exist", func() {
			err := configManager.SetActiveEnv(config.Env{
				Region: "non-existent",
			})
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("region"))

			err = configManager.SetActiveEnv(config.Env{
				Lifecycle: "non-existent",
			})
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("lifecycle"))
		})

		It("can create a new config file", func() {
			configFilePath := fmt.Sprintf("./test-config-%d.yaml", rand.Intn(1000))

			Expect(fileManager.SetConfigFile(configFilePath)).To(Succeed())
			Expect(fileManager.ConfigFile()).To(Equal(configFilePath))

			env := configManager.Config().ActiveEnv()
			Expect(env.Lifecycle).To(Equal("prod"))
			Expect(env.Region).To(Equal("na"))
			Expect(os.Remove(configFilePath)).NotTo(HaveOccurred())
		})

	})

	Context("using file config", func() {

		var (
			fileManager    *config.FileManager
			configManager  *config.ConfigManager
			configFilePath string
			input          *bytes.Buffer
			output         *bytes.Buffer
		)

		BeforeEach(func() {
			configFilePath = fmt.Sprintf("./test-config-%d.yaml", rand.Intn(1000000))
			Expect(config.WriteToFile(configFilePath, sampleConfig)).To(Succeed())

			input = bytes.NewBufferString("")
			output = bytes.NewBufferString("")
			formHandler := basic.NewFormHandler(input, output)
			fileManager = config.NewFileConfigManager("", formHandler)
			Expect(fileManager).NotTo(BeNil())
			Expect(fileManager.SetConfigFile(configFilePath)).To(Succeed())
			Expect(fileManager.ConfigFile()).To(Equal(configFilePath))

			configManager = config.NewConfigManager(fileManager, formHandler)
			Expect(configManager).NotTo(BeNil())

			config := configManager.Config()
			Expect(config.Region).To(Equal("na"))
			Expect(config.Lifecycle).To(Equal("prod"))
			ExpectConfigValue(configManager, "token").To(Equal("na-prod-token"))
			ExpectConfigValue(configManager, "username").To(Equal("na-prod-username"))
			ExpectConfigValue(configManager, "password").To(Equal("na-prod-password"))
		})

		AfterEach(func() {
			viper.Reset()
			Expect(os.Remove(configFilePath)).To(Succeed())
		})

		It("can get the config object", func() {
			config := fileManager.GetConfig()

			Expect(config).NotTo(BeNil())
			Expect(config.Envs).To(HaveLen(2))
			Expect(config.Envs["na-prod"].Region).To(Equal("na"))
			Expect(config.Envs["na-prod"].Lifecycle).To(Equal("prod"))
		})

		It("can set the active environment", func() {
			Expect(configManager.SetActiveEnv(config.Env{
				Region:    "na",
				Lifecycle: "stage",
			})).To(Succeed())

			env := configManager.Config().ActiveEnv()
			Expect(env.Region).To(Equal("na"))
			Expect(env.Lifecycle).To(Equal("stage"))
			ExpectConfigValue(configManager, "token").To(Equal("na-stage-token"))
			ExpectConfigValue(configManager, "username").To(Equal("na-stage-username"))
			ExpectConfigValue(configManager, "password").To(Equal("na-stage-password"))
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
			Expect(ok).To(BeTrue(), "user supplied value for api-key not found")
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

})

func ExpectConfigValue(configManager *config.ConfigManager, key string) Assertion {
	val, ok := configManager.GetConfigValue(key)
	Expect(ok).To(BeTrue(), "config value %s not found", key)
	return Expect(val)
}
