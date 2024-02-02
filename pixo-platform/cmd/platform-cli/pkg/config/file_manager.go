package config

import (
	"fmt"
	"github.com/PixoVR/pixo-golang-clients/pixo-platform/cmd/platform-cli/pkg/input"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
	"io"
	"os"
	"strings"
)

type fileManagerImpl struct {
	configFile string
	reader     io.Reader
	writer     io.Writer
}

func NewFileManager(cfgDir string) Manager {
	if cfgDir != "" {
		viper.AddConfigPath(cfgDir)
	}

	viper.SetConfigType("yaml")

	return &fileManagerImpl{}
}

func (f *fileManagerImpl) ReadConfigFile(configFile string) error {
	f.SetConfigFile(configFile)
	return f.readConfigFile()
}

func (f *fileManagerImpl) ConfigFile() string {
	return f.configFile
}

func (f *fileManagerImpl) SetConfigValue(key, value string) {
	configObj := f.GetConfig()
	activeEnv := configObj.ActiveEnv()
	activeEnv.Set(key, value)
	configObj.SetEnv(activeEnv)

	f.SetConfig(configObj)
}

func (f *fileManagerImpl) UnsetConfigValue(key, value string) {
	configObj := f.GetConfig()
	activeEnv := configObj.ActiveEnv()
	activeEnv.Unset(key)

	f.SetConfig(configObj)
}

func (f *fileManagerImpl) GetConfigValue(key string) (string, bool) {
	if val, ok := os.LookupEnv(f.formatEnvVar(key)); ok {
		return val, true
	}

	config := f.GetConfig()
	activeEnv := config.ActiveEnv()
	return activeEnv.Get(key)
}

func (f *fileManagerImpl) GetConfig() Config {
	var c Config
	_ = viper.Unmarshal(&c)
	if c.Envs == nil {
		c.Envs = make(map[string]Env)
	}

	return c
}

func (f *fileManagerImpl) SetConfig(configObj Config) {
	if configObj.Region != "" {
		viper.Set("region", configObj.Region)
	}

	if configObj.Lifecycle != "" {
		viper.Set("lifecycle", configObj.Lifecycle)
	}

	if configObj.Envs != nil {
		viper.Set("envs", configObj.Envs)
	}

	_ = viper.WriteConfig()
}

func (f *fileManagerImpl) readConfigFile() error {
	f.createConfigFileIfNotExists(f.configFile)
	viper.SetConfigFile(f.configFile)

	return viper.ReadInConfig()
}

func (f *fileManagerImpl) createConfigFileIfNotExists(configFile string) {
	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		if _, err = os.Create(configFile); err != nil {
			log.Error().Err(err).Msgf("error creating config file: %s", configFile)
		}
	}
}

func (f *fileManagerImpl) GetConfigValueOrAskUser(key, shortFlag string) string {
	if f.Reader() == nil {
		return ""
	}

	val, ok := f.GetConfigValue(key)
	if ok {
		return val
	}

	displayKey := strings.ReplaceAll(strings.ToUpper(key), "-", " ")
	if key == "password" {
		return input.ReadSensitiveFromUser(f.Writer(), displayKey)
	}

	return input.ReadFromUser(f.Reader(), f.Writer(), displayKey)
}

func (f *fileManagerImpl) SetActiveEnv(env Env) {
	if env.Region != "" {
		viper.Set("region", env.Region)
	}

	if env.Lifecycle != "" {
		viper.Set("lifecycle", env.Lifecycle)
	}

	_ = viper.WriteConfig()
}

func (f *fileManagerImpl) GetActiveEnv() Env {
	config := f.GetConfig()

	return config.ActiveEnv()
}

func (f *fileManagerImpl) formatEnvVar(key string) string {
	return fmt.Sprintf("PIXO_%s", strings.ReplaceAll(strings.ToUpper(key), "-", "_"))
}

func (f *fileManagerImpl) SetReader(r io.Reader) {
	f.reader = r
}

func (f *fileManagerImpl) Reader() io.Reader {
	return f.reader
}

func (f *fileManagerImpl) SetWriter(w io.Writer) {
	f.writer = w
}

func (f *fileManagerImpl) Writer() io.Writer {
	return f.writer
}

func (f *fileManagerImpl) SetConfigFile(configFile string) {
	f.configFile = configFile
}

func (f *fileManagerImpl) Lifecycle() string {
	if lifecycle, ok := os.LookupEnv("PIXO_LIFECYCLE"); ok {
		return lifecycle
	}

	return f.GetActiveEnv().Lifecycle
}

func (f *fileManagerImpl) Region() string {
	if region, ok := os.LookupEnv("PIXO_REGION"); ok {
		return region
	}

	return f.GetActiveEnv().Region
}
