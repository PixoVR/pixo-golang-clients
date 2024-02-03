package config

import (
	"github.com/spf13/viper"
	"os"
)

func (f *fileManagerImpl) ReadConfigFile(configFile string) error {
	f.SetConfigFile(configFile)
	return f.readConfigFile()
}

func (f *fileManagerImpl) ConfigFile() string {
	return f.configFile
}

func (f *fileManagerImpl) readConfigFile() error {
	f.createConfigFileIfNotExists(f.configFile)
	viper.SetConfigFile(f.configFile)

	return viper.ReadInConfig()
}

func (f *fileManagerImpl) createConfigFileIfNotExists(configFile string) {
	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		if _, err = os.Create(configFile); err != nil {
			_, _ = f.writer.Write([]byte("Error creating config file: " + err.Error()))
		}
	}
}
