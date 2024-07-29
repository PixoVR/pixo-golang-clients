package config_test

import (
	"bytes"
	"fmt"
	"github.com/PixoVR/pixo-golang-clients/pixo-platform/platform-cli/src/config"
	"github.com/PixoVR/pixo-golang-clients/pixo-platform/platform-cli/src/forms"
	"github.com/PixoVR/pixo-golang-clients/pixo-platform/platform-cli/src/forms/basic"
	"github.com/PixoVR/pixo-golang-clients/pixo-platform/platform-cli/src/printer"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"math/rand"
	"os"
)

var _ = Describe("Config Manager", func() {

	It("can get the active environment", func() {
		configObj := &config.Config{}
		env := configObj.ActiveEnv()
		Expect(configObj.Envs).To(HaveLen(0))
		Expect(env.Lifecycle).To(Equal("prod"))
		Expect(env.Region).To(Equal("na"))
	})

	Context("when the config file does not exist", func() {

		var (
			fileManager   *config.FileManager
			configManager *config.ConfigManager
		)

		BeforeEach(func() {
			fileManager = config.NewFileConfigManager("./test-config.yaml")
			Expect(fileManager).NotTo(BeNil())
			emojiPrinter := printer.NewEmojiPrinter(nil)
			configManager = config.NewConfigManager(fileManager, emojiPrinter)
			Expect(configManager).NotTo(BeNil())
		})

		AfterEach(func() {
			viper.Reset()
			Expect(os.Remove("./test-config.yaml")).To(Succeed())
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

			emojiPrinter := printer.NewEmojiPrinter(output)
			configManager = config.NewConfigManager(fileManager, emojiPrinter, formHandler)
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

		It("returns an error if the config value is not found", func() {
			input.WriteString("\n")

			val, ok := configManager.GetConfigValueOrAskUser("nonexistentconfig", nil)
			Expect(ok).To(BeFalse())
			Expect(val).To(BeEmpty())
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

		It("can use a form to get input values from the user if they're not supplied", func() {
			input.WriteString("some-val\n")
			configManager.SetConfigValue("config-key", "config-val")
			questions := []config.Value{
				{Question: forms.Question{Type: forms.Input, Key: "val", Prompt: "Enter val: "}},
				{Question: forms.Question{Type: forms.SensitiveInput, Key: "config-key", Prompt: "Enter config val: "}},
			}

			answers, err := configManager.GetValuesOrSubmitForm(questions, &cobra.Command{})

			Expect(err).NotTo(HaveOccurred())
			Expect(answers).To(HaveLen(len(questions)))

			Expect(forms.String(answers["val"])).To(Equal("some-val"))
			Expect(output.String()).To(ContainSubstring("Enter val: "))

			Expect(answers).To(HaveKeyWithValue("config-key", "config-val"))
			Expect(output.String()).NotTo(ContainSubstring("Enter config val: "))
		})

		It("can use a form to get confirm values from the user if they're not supplied", func() {
			input.WriteString("no\n")
			configManager.SetConfigValue("config-confirm", "yes")
			questions := []config.Value{
				{Question: forms.Question{Type: forms.Confirm, Key: "confirm", Prompt: "Enter confirm val: "}},
				{Question: forms.Question{Type: forms.Confirm, Key: "config-confirm", Prompt: "Enter config confirm val: "}},
			}

			answers, err := configManager.GetValuesOrSubmitForm(questions, &cobra.Command{})

			Expect(err).NotTo(HaveOccurred())
			Expect(answers).To(HaveLen(len(questions)))

			Expect(forms.Bool(answers["confirm"])).To(BeFalse())
			Expect(output.String()).To(ContainSubstring("Enter confirm val: "))

			Expect(forms.Bool(answers["config-confirm"])).To(BeTrue())
			Expect(output.String()).NotTo(ContainSubstring("Enter confirm config val: "))
		})

		It("can use a form to get select values from the user if they're not supplied", func() {
			input.WriteString("one\nthree\n")
			configManager.SetConfigValue("config-select", "two")
			configManager.SetConfigValue("config-select-id", "four")
			questions := []config.Value{
				{Question: forms.Question{
					Type:   forms.Select,
					Key:    "select",
					Prompt: "Enter select val: ",
					Options: []forms.Option{
						{Label: "one"},
						{Label: "two"},
					},
				}},
				{Question: forms.Question{
					Type:   forms.Select,
					Key:    "config-select",
					Prompt: "Enter config select vals: ",
					Options: []forms.Option{
						{Label: "one"},
						{Label: "two"},
					},
				}},
				{Question: forms.Question{
					Type:   forms.SelectID,
					Key:    "select-id",
					Prompt: "Enter select id val: ",
					Options: []forms.Option{
						{Label: "three", Value: "3"},
						{Label: "four", Value: "4"},
					},
				}},
				{Question: forms.Question{
					Type:   forms.SelectID,
					Key:    "config-select-id",
					Prompt: "Enter config select id vals: ",
					Options: []forms.Option{
						{Label: "three", Value: "3"},
						{Label: "four", Value: "4"},
					},
				}},
			}

			answers, err := configManager.GetValuesOrSubmitForm(questions, &cobra.Command{})

			Expect(err).NotTo(HaveOccurred())
			Expect(answers).To(HaveLen(len(questions)))

			Expect(forms.String(answers["select"])).To(Equal("one"))
			Expect(output.String()).To(ContainSubstring("Enter select val: "))

			Expect(forms.String(answers["config-select"])).To(Equal("two"))
			Expect(output.String()).NotTo(ContainSubstring("Enter config select vals: "))

			Expect(forms.Int(answers["select-id"])).To(Equal(3))
			Expect(output.String()).To(ContainSubstring("Enter select id val: "))

			Expect(forms.Int(answers["config-select-id"])).To(Equal(4))
			Expect(output.String()).NotTo(ContainSubstring("Enter config select id vals: "))
		})

		It("can use a form to get multiselect values from the user if they're not supplied", func() {
			input.WriteString("one,two\ntwo,one\n")
			configManager.SetConfigValue("config-multiselect", "three,four")
			configManager.SetConfigValue("config-multiselect-ids", "four,three")
			questions := []config.Value{
				{Question: forms.Question{
					Type:   forms.MultiSelect,
					Key:    "multiselect",
					Prompt: "Enter multiselect vals: ",
					Options: []forms.Option{
						{Label: "one"},
						{Label: "two"},
					},
				}},
				{Question: forms.Question{
					Type:   forms.MultiSelect,
					Key:    "config-multiselect",
					Prompt: "Enter config multiselect vals: ",
					Options: []forms.Option{
						{Label: "three"},
						{Label: "four"},
					},
				}},
				{Question: forms.Question{
					Type:   forms.MultiSelectIDs,
					Key:    "multiselect-ids",
					Prompt: "Enter multiselect id vals: ",
					Options: []forms.Option{
						{Label: "one", Value: "1"},
						{Label: "two", Value: "2"},
					},
				}},
				{Question: forms.Question{
					Type:   forms.MultiSelectIDs,
					Key:    "config-multiselect-ids",
					Prompt: "Enter config multiselect id vals: ",
					Options: []forms.Option{
						{Label: "three", Value: "3"},
						{Label: "four", Value: "4"},
					},
				}},
			}

			answers, err := configManager.GetValuesOrSubmitForm(questions, &cobra.Command{})

			Expect(err).NotTo(HaveOccurred())
			Expect(answers).To(HaveLen(len(questions)))

			Expect(forms.StringSlice(answers["multiselect"])).To(Equal([]string{"one", "two"}))
			Expect(output.String()).To(ContainSubstring("Enter multiselect vals: "))

			Expect(forms.StringSlice(answers["config-multiselect"])).To(Equal([]string{"three", "four"}))
			Expect(output.String()).NotTo(ContainSubstring("Enter config multiselect vals: "))

			Expect(forms.IntSlice(answers["multiselect-ids"])).To(Equal([]int{2, 1}))
			Expect(output.String()).To(ContainSubstring("Enter multiselect id vals: "))

			Expect(forms.IntSlice(answers["config-multiselect-ids"])).To(Equal([]int{4, 3}))
			Expect(output.String()).NotTo(ContainSubstring("Enter config multiselect id vals: "))
		})

	})

})
