package config

import (
	"github.com/spf13/viper"
	"strconv"
)

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

func (f *fileManagerImpl) SetActiveEnv(env Env) {
	if env.Region != "" {
		viper.Set("region", env.Region)
	}

	if env.Lifecycle != "" {
		viper.Set("lifecycle", env.Lifecycle)
	}

	_ = viper.WriteConfig()
}

func (f *fileManagerImpl) SetConfigValue(key, value string) {
	configObj := f.GetConfig()
	activeEnv := configObj.ActiveEnv()
	activeEnv.Set(key, value)
	configObj.SetEnv(activeEnv)

	f.SetConfig(configObj)
}

func (f *fileManagerImpl) SetIntConfigValue(key string, value int) {
	f.SetConfigValue(key, strconv.Itoa(value))
}

func (f *fileManagerImpl) SetBoolConfigValue(key string, value bool) {
	f.SetConfigValue(key, strconv.FormatBool(value))
}

func (f *fileManagerImpl) UnsetConfigValue(key string) {
	configObj := f.GetConfig()
	activeEnv := configObj.ActiveEnv()
	activeEnv.Unset(key)

	f.SetConfig(configObj)
}
