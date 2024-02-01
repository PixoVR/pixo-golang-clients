package config

import (
	"github.com/PixoVR/pixo-golang-clients/pixo-platform/cmd/platform-cli/pkg/input"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
	"os"
)

type fileManagerImpl struct {
	baseManagerImpl
}

func NewFileManager(cfgDir string) Manager {
	if cfgDir != "" {
		viper.AddConfigPath(cfgDir)
	}

	viper.SetConfigType("yaml")

	m := &fileManagerImpl{}
	m.SetActiveEnv(Env{Region: "na", Lifecycle: "prod"})

	return m
}

func (m *fileManagerImpl) ReadConfigFile(configFile string) error {
	m.SetConfigFile(configFile)
	return m.readConfigFile()
}

func (m *fileManagerImpl) ConfigFile() string {
	return m.configFile
}

func (m *fileManagerImpl) SetConfigValue(key, value string) error {
	viper.Set(key, value)
	return m.readConfigFile()
}

func (m *fileManagerImpl) GetConfigValue(key string) (string, bool) {
	config := m.GetConfig()
	activeEnv := config.ActiveEnv()
	return activeEnv.Get(key)
}

func (m *fileManagerImpl) GetConfig() Config {
	var c Config
	_ = viper.Unmarshal(&c)
	return c
}

func (m *fileManagerImpl) SetConfig(configObj Config) {
	viper.Set("region", configObj.Region)
	viper.Set("lifecycle", configObj.Lifecycle)
	//viper.Set("envs", configObj.Envs)
	//viper.Set("config", configObj)
	_ = viper.WriteConfig()
}

func (m *fileManagerImpl) readConfigFile() error {
	m.createConfigFileIfNotExists(m.configFile)
	viper.SetConfigFile(m.configFile)

	if err := viper.ReadInConfig(); err != nil {
		return err
	}
	return nil
}

func (m *fileManagerImpl) createConfigFileIfNotExists(configFile string) {
	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		if _, err = os.Create(configFile); err != nil {
			log.Error().Err(err).Msgf("error creating config file: %s", configFile)
		}
	}
}

func (m *fileManagerImpl) GetOrAsk(key string) string {
	if m.Reader() == nil {
		return ""
	}

	val, ok := m.GetConfigValue(key)
	if ok {
		return val
	}

	if key == "password" {
		return input.ReadSensitiveFromUser(m.Writer(), key)
	}

	return input.ReadFromUser(m.Reader(), m.Writer(), key)
}

func (m *fileManagerImpl) SetActiveEnv(env Env) {
	configObj := m.GetConfig()

	configObj.Region = env.Region
	configObj.Lifecycle = env.Lifecycle

	//configObj.SetEnv(env)

	m.SetConfig(configObj)
}

func (m *fileManagerImpl) GetActiveEnv() Env {
	config := m.GetConfig()

	return config.ActiveEnv()
}
