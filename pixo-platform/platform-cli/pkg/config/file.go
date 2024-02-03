package config

import (
	"github.com/spf13/viper"
	"os"
	"path/filepath"
	"strings"
)

func (f *fileManagerImpl) SetConfigFile(configFile string) error {
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

func (f *fileManagerImpl) ConfigFile() string {
	return f.configFile
}

func (f *fileManagerImpl) createConfigFileIfNotExists(configFile string) {
	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		if _, err = os.Create(configFile); err != nil {
			_, _ = f.writer.Write([]byte("Error creating config file: " + err.Error()))
		}
	}
}
