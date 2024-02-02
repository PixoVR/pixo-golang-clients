package config_test

import (
	"bytes"
	"fmt"
	"github.com/PixoVR/pixo-golang-clients/pixo-platform/cmd/platform-cli/pkg/config"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/spf13/viper"
	"math/rand"
	"os"
)

var _ = Describe("Config", func() {

	AfterEach(func() {
		viper.Reset()
	})

	It("can get the active environment", func() {
		configObj := config.Config{}
		env := configObj.ActiveEnv()
		Expect(configObj.Envs).To(HaveLen(0))
		Expect(env.Lifecycle).To(Equal("prod"))
		Expect(env.Region).To(Equal("na"))
	})

	Context("when the config file does not exist", func() {

		It("can initialize the config manager when setting the active environment", func() {
			configManager := config.NewFileManager("")
			Expect(configManager).NotTo(BeNil())
			env := config.Env{
				Region:    "na",
				Lifecycle: "stage",
			}

			configManager.SetActiveEnv(env)

			Expect(configManager.Region()).To(Equal(env.Region))
			Expect(configManager.Lifecycle()).To(Equal(env.Lifecycle))
		})

		It("can initialize the config manager when setting a value", func() {
			configManager := config.NewFileManager("")
			Expect(configManager).NotTo(BeNil())
			apiKey := "some-api-key"
			configManager.SetConfigValue("api-key", apiKey)
			ExpectConfigValueToEqual(configManager, "api-key", apiKey)
		})

		It("can create a new config file", func() {
			configManager := config.NewFileManager("")
			Expect(configManager).NotTo(BeNil())
			configFilePath := fmt.Sprintf("./test-config-%d.yaml", rand.Intn(1000))

			err := configManager.ReadConfigFile(configFilePath)
			Expect(err).NotTo(HaveOccurred())
			Expect(configManager.ConfigFile()).To(Equal(configFilePath))
			env := configManager.GetActiveEnv()
			Expect(env.Lifecycle).To(Equal("prod"))
			Expect(env.Region).To(Equal("na"))

			err = os.Remove(configFilePath)
		})

	})

	var (
		sampleConfig = config.Config{
			Lifecycle: "prod",
			Region:    "na",
			Envs: map[string]config.Env{
				"na-prod": {
					Region:    "na",
					Lifecycle: "prod",
					EnvMap: map[string]string{
						"api-key":  "na-prod-api-key",
						"token":    "na-prod-token",
						"username": "na-prod-username",
						"password": "na-prod-password",
					},
				},
				"na-stage": {
					Region:    "na",
					Lifecycle: "stage",
					EnvMap: map[string]string{
						"api-key":  "na-stage-api-key",
						"token":    "na-stage-token",
						"username": "na-stage-username",
						"password": "na-stage-password",
					},
				},
			},
		}
	)

	Context("using file config", Ordered, func() {

		var (
			fileConfigManager config.Manager
			configFilePath    string
		)

		BeforeEach(func() {
			configFilePath = fmt.Sprintf("./test-config-%d.yaml", rand.Intn(1000))
			err := config.WriteToFile(configFilePath, sampleConfig)
			Expect(err).NotTo(HaveOccurred())

			fileConfigManager = config.NewFileManager("")
			Expect(fileConfigManager).NotTo(BeNil())

			err = fileConfigManager.ReadConfigFile(configFilePath)
			Expect(err).NotTo(HaveOccurred())
			Expect(fileConfigManager.ConfigFile()).To(Equal(configFilePath))

			fileConfig := fileConfigManager.GetConfig()
			Expect(fileConfig.Region).To(Equal("na"))
			Expect(fileConfig.Lifecycle).To(Equal("prod"))
			ExpectConfigValueToEqual(fileConfigManager, "token", "na-prod-token")
			ExpectConfigValueToEqual(fileConfigManager, "username", "na-prod-username")
			ExpectConfigValueToEqual(fileConfigManager, "password", "na-prod-password")
		})

		AfterEach(func() {
			viper.Reset()
			err := os.Remove(configFilePath)
			Expect(err).NotTo(HaveOccurred())
		})

		It("can get the config object", func() {
			configObj := fileConfigManager.GetConfig()

			Expect(configObj).NotTo(BeNil())
			Expect(configObj.Envs).To(HaveLen(2))
			Expect(configObj.Envs["na-prod"].Region).To(Equal("na"))
			Expect(configObj.Envs["na-prod"].Lifecycle).To(Equal("prod"))
		})

		It("can set the active environment", func() {
			fileConfigManager.SetActiveEnv(config.Env{
				Region:    "na",
				Lifecycle: "stage",
			})

			env := fileConfigManager.GetActiveEnv()
			Expect(env.Region).To(Equal("na"))
			Expect(env.Lifecycle).To(Equal("stage"))
			ExpectConfigValueToEqual(fileConfigManager, "token", "na-stage-token")
			ExpectConfigValueToEqual(fileConfigManager, "username", "na-stage-username")
			ExpectConfigValueToEqual(fileConfigManager, "password", "na-stage-password")
		})

		It("can set values in the active environment", func() {
			fileConfigManager.SetConfigValue("api-key", "new-prod-api-key")
			ExpectConfigValueToEqual(fileConfigManager, "api-key", "new-prod-api-key")
		})

		It("can prioritize environment variables", func() {
			Expect(os.Setenv("PIXO_REGION", "saudi")).NotTo(HaveOccurred())
			Expect(os.Setenv("PIXO_LIFECYCLE", "dev")).NotTo(HaveOccurred())
			Expect(os.Setenv("PIXO_API_KEY", "saudi-dev-api-key")).NotTo(HaveOccurred())

			Expect(fileConfigManager.Region()).To(Equal("saudi"))
			Expect(fileConfigManager.Lifecycle()).To(Equal("dev"))
			ExpectConfigValueToEqual(fileConfigManager, "api-key", "saudi-dev-api-key")

			Expect(os.Unsetenv("PIXO_REGION")).NotTo(HaveOccurred())
			Expect(os.Unsetenv("PIXO_LIFECYCLE")).NotTo(HaveOccurred())
			Expect(os.Unsetenv("PIXO_API_KEY")).NotTo(HaveOccurred())
		})

		It("can get a config value if it exists instead of asking the user", func() {
			username := "new-username"
			inputReader := bytes.NewBufferString(username + "\n")
			fileConfigManager.SetReader(inputReader)
			outputWriter := bytes.NewBufferString("")
			fileConfigManager.SetWriter(outputWriter)

			val := fileConfigManager.GetConfigValueOrAskUser("username", "u")
			Expect(outputWriter.String()).To(BeEmpty())
			Expect(val).To(Equal("na-prod-username"))
		})

		It("can ask the user if the config doesnt exist", func() {
			fileConfigManager.UnsetConfigValue("api-key", "")
			apiKey := "new-api-key"
			inputReader := bytes.NewBufferString(apiKey + "\n")
			fileConfigManager.SetReader(inputReader)
			outputWriter := bytes.NewBufferString("")
			fileConfigManager.SetWriter(outputWriter)

			val := fileConfigManager.GetConfigValueOrAskUser("api-key", "a")
			Expect(outputWriter.String()).To(ContainSubstring("Enter API KEY: "))
			Expect(val).To(Equal(apiKey))
		})

	})

})

func ExpectConfigValueToEqual(configManager config.Manager, key, expectedValue string) {
	val, ok := configManager.GetConfigValue(key)
	Expect(val).To(Equal(expectedValue))
	Expect(ok).To(BeTrue())
}
