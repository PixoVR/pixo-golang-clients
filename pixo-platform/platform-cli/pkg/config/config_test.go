package config_test

import (
	"bytes"
	"fmt"
	"github.com/PixoVR/pixo-golang-clients/pixo-platform/platform-cli/pkg/config"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/spf13/cobra"
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

		var (
			configManager config.Manager
		)

		BeforeEach(func() {
			configManager = config.NewFileManager("./test-config.yaml")
			Expect(configManager).NotTo(BeNil())
		})

		AfterEach(func() {
			viper.Reset()
			err := os.Remove("./test-config.yaml")
			Expect(err).NotTo(HaveOccurred())
		})

		It("can initialize the config manager when setting the active environment", func() {
			env := config.Env{
				Region:    "na",
				Lifecycle: "stage",
			}

			Expect(configManager.SetActiveEnv(env)).To(Succeed())

			Expect(configManager.Region()).To(Equal(env.Region))
			Expect(configManager.Lifecycle()).To(Equal(env.Lifecycle))
		})

		It("can use stdout and stdin by default", func() {
			n, err := configManager.Write([]byte("hello"))
			Expect(n).To(Equal(5))
			Expect(err).NotTo(HaveOccurred())
			n, err = configManager.Read([]byte{})
			Expect(n).To(Equal(0))
			Expect(err).NotTo(HaveOccurred())
		})

		It("can output a messages with emojis", func() {
			outputWriter := bytes.NewBufferString("")
			configManager.SetWriter(outputWriter)
			msg := "hello world\n"

			configManager.Print(msg)
			Expect(outputWriter.String()).To(Equal(msg))
			outputWriter.Reset()

			configManager.Println(msg)
			Expect(outputWriter.String()).To(Equal(msg + "\n"))
			outputWriter.Reset()

			msg = fmt.Sprintf(":rocket:hello %s", "world")
			expectedMsg := "🚀 hello world"
			configManager.Printf(msg)
			Expect(outputWriter.String()).To(Equal(expectedMsg))
		})

		It("can initialize the config manager when setting a value", func() {
			apiKey := "some-api-key"

			configManager.SetConfigValue("api-key", apiKey)

			ExpectConfigValueToEqual(configManager, "api-key", apiKey)
		})

		It("can return an error if the region does not exist", func() {
			err := configManager.SetActiveEnv(config.Env{
				Region: "non-existent",
			})
			Expect(err).To(HaveOccurred())

			err = configManager.SetActiveEnv(config.Env{
				Lifecycle: "non-existent",
			})
			Expect(err).To(HaveOccurred())
		})

		It("can create a new config file", func() {
			configFilePath := fmt.Sprintf("./test-config-%d.yaml", rand.Intn(1000))

			err := configManager.SetConfigFile(configFilePath)

			Expect(err).NotTo(HaveOccurred())
			Expect(configManager.ConfigFile()).To(Equal(configFilePath))
			env := configManager.GetActiveEnv()
			Expect(env.Lifecycle).To(Equal("prod"))
			Expect(env.Region).To(Equal("na"))
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

			err = fileConfigManager.SetConfigFile(configFilePath)
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
			Expect(fileConfigManager.SetActiveEnv(config.Env{
				Region:    "na",
				Lifecycle: "stage",
			})).To(Succeed())

			env := fileConfigManager.GetActiveEnv()
			Expect(env.Region).To(Equal("na"))
			Expect(env.Lifecycle).To(Equal("stage"))
			ExpectConfigValueToEqual(fileConfigManager, "token", "na-stage-token")
			ExpectConfigValueToEqual(fileConfigManager, "username", "na-stage-username")
			ExpectConfigValueToEqual(fileConfigManager, "password", "na-stage-password")
		})

		It("can use the same local config regardless of region", func() {
			Expect(fileConfigManager.SetActiveEnv(config.Env{
				Region:    "na",
				Lifecycle: "local",
			})).To(Succeed())

			fileConfigManager.SetConfigValue("api-key", "local-api-key")

			Expect(fileConfigManager.SetActiveEnv(config.Env{
				Region: "saudi",
			})).To(Succeed())
			ExpectConfigValueToEqual(fileConfigManager, "api-key", "local-api-key")
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

			val, ok := fileConfigManager.GetConfigValueOrAskUser("username", nil)
			Expect(ok).To(BeTrue())
			Expect(val).To(Equal("na-prod-username"))
			Expect(outputWriter.String()).To(BeEmpty())
		})

		It("can ask the user if the config doesnt exist", func() {
			fileConfigManager.UnsetConfigValue("api-key")
			apiKey := "new-api-key"
			inputReader := bytes.NewBufferString(apiKey + "\n")
			fileConfigManager.SetReader(inputReader)
			outputWriter := bytes.NewBufferString("")
			fileConfigManager.SetWriter(outputWriter)

			val, ok := fileConfigManager.GetConfigValueOrAskUser("api-key", nil)
			Expect(ok).To(BeTrue())
			Expect(val).To(Equal(apiKey))
			Expect(outputWriter.String()).To(ContainSubstring("Enter API KEY: "))
		})

		It("can ask the user if the config doesnt exist for an int value", func() {
			fileConfigManager.SetIntConfigValue("val", 1)
			val, ok := fileConfigManager.GetIntConfigValueOrAskUser("val", nil)
			Expect(ok).To(BeTrue())
			Expect(val).To(Equal(1))
		})

		It("can set an int config value and get it", func() {
			fileConfigManager.SetIntConfigValue("port", 8080)
			val, ok := fileConfigManager.GetIntConfigValue("port")
			Expect(ok).To(BeTrue())
			Expect(val).To(Equal(8080))
		})

		It("can set a bool config value and get it", func() {
			fileConfigManager.SetBoolConfigValue("is-active", true)
			val, ok := fileConfigManager.GetBoolConfigValue("is-active")
			Expect(ok).To(BeTrue())
			Expect(val).To(Equal(true))
		})

		It("can prioritize the value from a flag over everything else", func() {
			cmd := &cobra.Command{}
			cmd.Flags().String("api-key", "flag", "api key")
			_ = cmd.Flags().Set("api-key", "flag-api-key")

			val, ok := fileConfigManager.GetFlagOrConfigValue("api-key", cmd)
			Expect(ok).To(BeTrue())
			Expect(val).To(Equal("flag-api-key"))
		})

		It("can get an int value from a flag", func() {
			cmd := &cobra.Command{}
			cmd.Flags().Int("port", 8080, "port")
			_ = cmd.Flags().Set("port", "9090")

			val, ok := fileConfigManager.GetIntFlagOrConfigValue("port", cmd)
			Expect(ok).To(BeTrue())
			Expect(val).To(Equal(9090))
		})

		It("can get a bool value from a flag", func() {
			cmd := &cobra.Command{}
			cmd.Flags().Bool("is-active", true, "is active")
			_ = cmd.Flags().Set("is-active", "false")

			val, ok := fileConfigManager.GetBoolFlagOrConfigValue("is-active", cmd)
			Expect(ok).To(BeTrue())
			Expect(val).To(Equal(false))
		})

		It("can read from the user", func() {
			inputReader := bytes.NewBufferString("new-username\n")
			outputWriter := bytes.NewBufferString("")
			fileConfigManager.SetReader(inputReader)
			fileConfigManager.SetWriter(outputWriter)

			val := fileConfigManager.ReadFromUser("username")
			Expect(val).To(Equal("new-username"))
			Expect(outputWriter.String()).To(ContainSubstring("Enter username: "))
		})

		It("can read a sensitive value from the user", func() {
			inputReader := bytes.NewBufferString("new-password\n")
			outputWriter := bytes.NewBufferString("")
			fileConfigManager.SetReader(inputReader)
			fileConfigManager.SetWriter(outputWriter)

			val := fileConfigManager.ReadSensitiveFromUser("password")
			Expect(val).To(Equal("new-password"))
			Expect(outputWriter.String()).To(ContainSubstring("Enter password: "))
			Expect(outputWriter.String()).NotTo(ContainSubstring("new-password"))
		})

	})

})

func ExpectConfigValueToEqual(configManager config.Manager, key, expectedValue string) {
	val, ok := configManager.GetConfigValue(key)
	Expect(val).To(Equal(expectedValue))
	Expect(ok).To(BeTrue())
}
