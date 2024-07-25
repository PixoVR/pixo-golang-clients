package config

import (
	"github.com/PixoVR/pixo-golang-clients/pixo-platform/platform-cli/src/forms"
	"github.com/PixoVR/pixo-golang-clients/pixo-platform/platform-cli/src/forms/basic"
	"github.com/spf13/cobra"
	"io"
	"os"
	"strconv"
	"strings"
)

type ConfigManager struct {
	formHandler forms.FormHandler
	manager     Manager
}

func NewConfigManager(manager Manager, formHandlers ...forms.FormHandler) *ConfigManager {
	var formHandler forms.FormHandler
	if len(formHandlers) > 0 {
		formHandler = formHandlers[0]
	} else {
		formHandler = basic.NewFormHandler(os.Stdin, os.Stdout)
	}

	return &ConfigManager{
		formHandler: formHandler,
		manager:     manager,
	}
}

func (c *ConfigManager) ActiveEnv() Env {
	return c.manager.GetConfig().ActiveEnv()
}

func (c *ConfigManager) Region() string {
	envVar, ok := lookupConfigEnv("region")
	if ok {
		return envVar
	}
	return c.ActiveEnv().Region
}

func (c *ConfigManager) Lifecycle() string {
	envVar, ok := lookupConfigEnv("lifecycle")
	if ok {
		return envVar
	}
	return c.ActiveEnv().Lifecycle
}

func (c *ConfigManager) SetActiveEnv(env Env) error {
	if err := env.Validate(); err != nil {
		return err
	}

	config := c.manager.GetConfig()
	if config == nil {
		config = &Config{}
	}

	if env.Region != "" {
		config.Region = env.Region
	}

	if env.Lifecycle != "" {
		config.Lifecycle = env.Lifecycle
	}

	c.manager.SetConfig(*config)
	return nil
}

func (c *ConfigManager) Config() *Config {
	config := c.manager.GetConfig()

	if config == nil {
		config = &Config{}
	}

	if config.Lifecycle == "" {
		config.Lifecycle = "prod"
	}

	if config.Region == "" {
		config.Region = "na"
	}

	return config
}

func (c *ConfigManager) GetConfigValue(key string) (string, bool) {
	if val, ok := os.LookupEnv(formatEnvVarName(key)); ok {
		return val, true
	}

	activeEnv := c.Config().ActiveEnv()
	return activeEnv.Get(key)
}

func (c *ConfigManager) SetConfigValue(key, value string) {
	config := c.Config()
	activeEnv := config.ActiveEnv()
	activeEnv.Set(key, value)
	config.SetEnv(activeEnv)

	c.manager.SetConfig(*config)
}

func (c *ConfigManager) UnsetConfigValue(key string) {
	config := c.Config()
	activeEnv := config.ActiveEnv()
	activeEnv.Unset(key)
	config.SetEnv(activeEnv)

	c.manager.SetConfig(*config)
}

func (c *ConfigManager) GetIntConfigValue(key string) (int, bool) {
	val, ok := c.GetConfigValue(key)
	if !ok {
		return 0, false
	}

	return ToInt(val)
}

func (c *ConfigManager) SetIntConfigValue(key string, value int) {
	c.SetConfigValue(key, strconv.Itoa(value))
}

func (c *ConfigManager) GetBoolConfigValue(key string) (bool, bool) {
	val, ok := c.GetConfigValue(key)
	if !ok {
		return false, false
	}

	return ToBool(val)
}

func (c *ConfigManager) SetBoolConfigValue(key string, value bool) {
	c.SetConfigValue(key, strconv.FormatBool(value))
}

func (c *ConfigManager) GetConfigValueOrAskUser(key string, cmd *cobra.Command) (string, bool) {
	flag, ok := c.GetFlagValue(key, cmd)
	if ok {
		return flag, true
	}

	val, ok := c.GetConfigValue(key)
	if ok {
		return val, true
	}

	displayKey := strings.ReplaceAll(strings.ToUpper(key), "-", " ")
	if strings.ToLower(key) == "password" {
		response, err := c.formHandler.GetSensitiveResponseFromUser(displayKey)
		if err != nil {
			return "", false
		}
		return response, response != ""
	}

	val, err := c.formHandler.GetResponseFromUser(displayKey)
	return val, err == nil && val != ""
}

func (c *ConfigManager) GetIntConfigValueOrAskUser(key string, cmd *cobra.Command) (int, bool) {
	flag, ok := c.GetIntFlagValue(key, cmd)
	if ok {
		return flag, true
	}

	val, ok := c.GetIntConfigValue(key)
	if ok {
		return val, true
	}

	displayKey := strings.ReplaceAll(strings.ToUpper(key), "-", " ")

	strVal, err := c.formHandler.GetResponseFromUser(displayKey)
	if err != nil {
		return 0, false
	}

	return ToInt(strVal)
}

func (c *ConfigManager) GetBoolConfigValueOrAskUser(key string, cmd *cobra.Command) (bool, bool) {
	flag, ok := c.GetBoolFlagValue(key, cmd)
	if ok {
		return flag, true
	}

	val, ok := c.GetBoolConfigValue(key)
	if ok {
		return val, true
	}

	displayKey := strings.ReplaceAll(strings.ToUpper(key), "-", " ")
	strVal, err := c.formHandler.GetResponseFromUser(displayKey)
	if err != nil {
		return false, false
	}

	boolVal, err := strconv.ParseBool(strVal)
	return boolVal, err == nil
}

func (c *ConfigManager) GetFlagOrConfigValue(key string, cmd *cobra.Command) (string, bool) {
	val, ok := c.GetFlagValue(key, cmd)
	if ok {
		return val, true
	}

	return c.GetConfigValue(key)
}

func (c *ConfigManager) GetIntFlagOrConfigValue(key string, cmd *cobra.Command) (int, bool) {
	val, ok := c.GetFlagOrConfigValue(key, cmd)
	if !ok {
		return 0, false
	}

	return ToInt(val)
}

func (c *ConfigManager) GetBoolFlagOrConfigValue(key string, cmd *cobra.Command) (bool, bool) {
	val, ok := c.GetFlagOrConfigValue(key, cmd)
	if !ok {
		return false, false
	}

	return ToBool(val)
}

func (c *ConfigManager) GetFlagValue(key string, cmd *cobra.Command) (string, bool) {
	if cmd == nil {
		return "", false
	}

	flag := cmd.Flag(key)
	if flag != nil {
		if val := flag.Value.String(); val != "" {
			return val, true
		}
	}

	return "", false
}

func (c *ConfigManager) GetIntFlagValue(key string, cmd *cobra.Command) (int, bool) {
	val, ok := c.GetFlagValue(key, cmd)
	if !ok {
		return 0, false
	}

	return ToInt(val)
}

func (c *ConfigManager) GetBoolFlagValue(key string, cmd *cobra.Command) (bool, bool) {
	val, ok := c.GetFlagValue(key, cmd)
	if !ok {
		return false, false
	}

	return ToBool(val)
}

func (c *ConfigManager) GetSensitiveConfigValueOrAskUser(key string, cmd *cobra.Command) (string, bool) {
	val, ok := c.GetConfigValue(key)
	if ok {
		return val, true
	}

	val, err := c.formHandler.GetSensitiveResponseFromUser(key)
	return val, err == nil
}

func (c *ConfigManager) GetFlagOrConfigValueOrAskUser(key string, cmd *cobra.Command) (string, bool) {
	val, ok := c.GetFlagOrConfigValue(key, cmd)
	if ok {
		return val, true
	}

	val, err := c.formHandler.GetResponseFromUser(key)
	return val, err == nil
}

func (c *ConfigManager) GetSensitiveFlagOrConfigValueOrAskUser(key string, cmd *cobra.Command) (string, bool) {
	val, ok := c.GetFlagOrConfigValue(key, cmd)
	if ok {
		return val, true
	}

	val, err := c.formHandler.GetSensitiveResponseFromUser(key)
	return val, err == nil
}

func (c *ConfigManager) GetIntFlagOrConfigValueOrAskUser(key string, cmd *cobra.Command) (int, bool) {
	val, ok := c.GetFlagOrConfigValue(key, cmd)
	if !ok {
		var err error
		val, err = c.formHandler.GetResponseFromUser(key)
		if err != nil {
			return 0, false
		}
	}

	return ToInt(val)
}

func (c *ConfigManager) GetBoolFlagOrConfigValueOrAskUser(key string, cmd *cobra.Command) (bool, bool) {
	val, ok := c.GetFlagOrConfigValue(key, cmd)
	if !ok {
		var err error
		val, err = c.formHandler.GetResponseFromUser(key)
		if err != nil {
			return false, false
		}
	}

	return ToBool(val)
}

func (c *ConfigManager) SetReader(reader io.Reader) {
	c.formHandler.SetReader(reader)
}

func (c *ConfigManager) SetWriter(writer io.Writer) {
	c.formHandler.SetWriter(writer)
}
