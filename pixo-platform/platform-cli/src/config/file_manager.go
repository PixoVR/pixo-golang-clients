package config

import (
	"github.com/PixoVR/pixo-golang-clients/pixo-platform/platform-cli/src/forms"
	"github.com/PixoVR/pixo-golang-clients/pixo-platform/platform-cli/src/forms/basic"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
	"os"
	"path/filepath"
	"strings"
)

var _ Manager = (*FileManager)(nil)

type FileManager struct {
	config      *Config
	configFile  string
	formHandler forms.FormHandler
}

func NewFileConfigManager(cfgFile string, formHandlers ...forms.FormHandler) *FileManager {
	var formHandler forms.FormHandler
	if len(formHandlers) > 0 {
		formHandler = formHandlers[0]
	} else {
		formHandler = basic.NewFormHandler(os.Stdin, os.Stdout)
	}

	m := &FileManager{
		config:      &Config{},
		formHandler: formHandler,
	}

	if cfgFile != "" {
		_ = m.SetConfigFile(cfgFile)
	}

	return m
}

func (f *FileManager) GetConfig() *Config {
	c := &Config{}
	_ = viper.Unmarshal(c)
	if c.Envs == nil {
		c.Envs = make(map[string]Env)
	}
	return c
}

func (f *FileManager) SetConfig(config Config) {
	if config.Region != "" {
		viper.Set("region", config.Region)
	}

	if config.Lifecycle != "" {
		viper.Set("lifecycle", config.Lifecycle)
	}

	if config.Envs != nil {
		viper.Set("envs", config.Envs)
	}

	_ = viper.WriteConfig()
}

func (f *FileManager) SetConfigFile(configFile string) error {
	f.configFile = configFile
	f.createConfigFileIfNotExists(f.configFile)
	viper.SetConfigFile(f.configFile)

	configFileExt := strings.ReplaceAll(filepath.Ext(configFile), ".", "")
	if configFileExt == "" {
		configFileExt = "yaml"
	}

	viper.SetConfigType(configFileExt)

	return viper.ReadInConfig()
}

func (f *FileManager) ConfigFile() string {
	return f.configFile
}

func (f *FileManager) createConfigFileIfNotExists(configFile string) {
	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		if _, err = os.Create(configFile); err != nil {
			log.Error().Err(err).Msgf("Failed to create config file %s", configFile)
		}
	}
}

func WriteToFile(path string, configFile Config) error {
	configFileBytes, err := yaml.Marshal(configFile)
	if err != nil {
		return err
	}

	if err = os.WriteFile(path, configFileBytes, 0644); err != nil {
		return err
	}

	return nil
}
