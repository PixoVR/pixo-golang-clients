package config_test

import (
	"bytes"
	"fmt"
	"github.com/PixoVR/pixo-golang-clients/pixo-platform/cmd/platform-cli/pkg/config"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"math/rand"
	"os"
)

var _ = Describe("GetConfig", func() {

	Context("when the config file does not exist", func() {

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
						"apiKey":   "na-prod-api-key",
						"token":    "na-prod-token",
						"username": "na-prod-username",
						"password": "na-prod-password",
					},
				},
				"na-stage": {
					Region:    "na",
					Lifecycle: "stage",
					EnvMap: map[string]string{
						"apiKey":   "na-stage-api-key",
						"token":    "na-stage-token",
						"username": "na-stage-username",
						"password": "na-stage-password",
					},
				},
			},
		}
	)

	Context("using in memory config", Ordered, func() {

		var (
			memoryConfigManager config.Manager
		)

		BeforeAll(func() {
			memoryConfigManager = config.NewInMemoryManager()
			Expect(memoryConfigManager).NotTo(BeNil())
			memoryConfigManager.SetConfig(sampleConfig)
			Expect(memoryConfigManager.Lifecycle()).To(Equal("prod"))
			Expect(memoryConfigManager.Region()).To(Equal("na"))
			username, ok := memoryConfigManager.GetConfigValue("username")
			Expect(ok).To(BeTrue())
			Expect(username).To(Equal("na-prod-username"))
		})

		It("can set the active environment", func() {
			memoryConfigManager.SetActiveEnv(config.Env{
				Region:    "na",
				Lifecycle: "stage",
			})

			Expect(memoryConfigManager.Region()).To(Equal("na"))
			Expect(memoryConfigManager.Lifecycle()).To(Equal("stage"))
			Expect(memoryConfigManager.GetActiveEnv().Region).To(Equal("na"))
			Expect(memoryConfigManager.GetActiveEnv().Lifecycle).To(Equal("stage"))
			ExpectConfigValueToEqual(memoryConfigManager, "token", "na-stage-token")
			ExpectConfigValueToEqual(memoryConfigManager, "username", "na-stage-username")
			ExpectConfigValueToEqual(memoryConfigManager, "password", "na-stage-password")
		})

	})

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
			fileConfig := fileConfigManager.GetConfig()
			Expect(fileConfig.Region).To(Equal("na"))
			Expect(fileConfig.Lifecycle).To(Equal("prod"))
			Expect(fileConfigManager.ConfigFile()).To(Equal(configFilePath))
			ExpectConfigValueToEqual(fileConfigManager, "token", "na-prod-token")
			ExpectConfigValueToEqual(fileConfigManager, "username", "na-prod-username")
			ExpectConfigValueToEqual(fileConfigManager, "password", "na-prod-password")
		})

		AfterEach(func() {
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

		It("can set values the active environment", func() {
			fileConfigManager.SetActiveEnv(config.Env{
				Region:    "na",
				Lifecycle: "stage",
			})

			Expect(fileConfigManager.SetConfigValue("api-key", "new-api-key")).To(Succeed())
			Expect(fileConfigManager.GetConfigValue("api-key")).To(Equal("new-api-key"))
		})

		It("can prioritize environment variables", func() {
			Expect(os.Setenv("PIXO_API_KEY", "saudi-dev-api-key")).NotTo(HaveOccurred())
			Expect(os.Setenv("PIXO_REGION", "saudi")).NotTo(HaveOccurred())
			Expect(os.Setenv("PIXO_LIFECYCLE", "dev")).NotTo(HaveOccurred())
			Expect(os.Setenv("PIXO_TOKEN", "saudi-dev-token")).NotTo(HaveOccurred())
			Expect(os.Setenv("PIXO_USERNAME", "saudi-dev-username")).NotTo(HaveOccurred())
			Expect(os.Setenv("PIXO_PASSWORD", "saudi-dev-password")).NotTo(HaveOccurred())

			Expect(fileConfigManager.GetConfigValue("api-key")).To(Equal("saudi-dev-api-key"))
			Expect(fileConfigManager.GetConfigValue("region")).To(Equal("saudi"))
			Expect(fileConfigManager.GetConfigValue("lifecycle")).To(Equal("dev"))
			Expect(fileConfigManager.GetConfigValue("token")).To(Equal("saudi-dev-token"))
			Expect(fileConfigManager.GetConfigValue("username")).To(Equal("saudi-dev-username"))
			Expect(fileConfigManager.GetConfigValue("password")).To(Equal("saudi-dev-password"))

		})

		It("can get a config value if it exists instead of asking the user", func() {
			username := "fake-username"
			err := fileConfigManager.SetConfigValue("username", username)
			Expect(err).NotTo(HaveOccurred())
			inputReader := bytes.NewBufferString("new-api-key\n")
			fileConfigManager.SetReader(inputReader)
			outputWriter := bytes.NewBufferString("")
			fileConfigManager.SetWriter(outputWriter)

			val := fileConfigManager.GetOrAsk("username")
			Expect(outputWriter.String()).To(BeEmpty())
			Expect(val).To(Equal(username))
		})

		It("can ask the user if the config doesnt exist", func() {
			inputReader := bytes.NewBufferString("new-api-key\n")
			fileConfigManager.SetReader(inputReader)
			outputWriter := bytes.NewBufferString("")
			fileConfigManager.SetWriter(outputWriter)

			val := fileConfigManager.GetOrAsk("api-key")
			Expect(outputWriter.String()).To(ContainSubstring("Enter API Key: "))
			Expect(val).To(Equal("new-api-key"))
		})

	})

})

func ExpectConfigValueToEqual(configManager config.Manager, key, expectedValue string) {
	val, ok := configManager.GetConfigValue(key)
	Expect(val).To(Equal(expectedValue))
	Expect(ok).To(BeTrue())
}
