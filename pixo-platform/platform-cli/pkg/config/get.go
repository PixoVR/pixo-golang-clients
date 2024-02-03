package config

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
	"strconv"
	"strings"
)

func (f *fileManagerImpl) GetConfig() Config {
	var c Config
	_ = viper.Unmarshal(&c)
	if c.Envs == nil {
		c.Envs = make(map[string]Env)
	}

	return c
}

func (f *fileManagerImpl) GetActiveEnv() Env {
	config := f.GetConfig()

	return config.ActiveEnv()
}

func (f *fileManagerImpl) GetConfigValue(key string) (string, bool) {
	if val, ok := os.LookupEnv(f.formatEnvVar(key)); ok {
		return val, true
	}

	config := f.GetConfig()
	activeEnv := config.ActiveEnv()
	return activeEnv.Get(key)
}

func (f *fileManagerImpl) GetIntConfigValue(key string) (int, bool) {
	val, ok := f.GetConfigValue(key)
	if !ok {
		return 0, false
	}

	return ToInt(val)
}

func (f *fileManagerImpl) GetBoolConfigValue(key string) (bool, bool) {
	val, ok := f.GetConfigValue(key)
	if !ok {
		return false, false
	}

	return ToBool(val)
}

func (f *fileManagerImpl) GetConfigValueOrAskUser(key string, cmd *cobra.Command) (string, bool) {
	flag, ok := f.GetFlagValue(key, cmd)
	if ok {
		return flag, true
	}

	val, ok := f.GetConfigValue(key)
	if ok {
		return val, true
	}

	displayKey := strings.ReplaceAll(strings.ToUpper(key), "-", " ")
	if strings.ToLower(key) == "password" {
		return f.ReadSensitiveFromUser(displayKey), true
	}

	val = f.ReadFromUser(displayKey)
	return val, val != ""
}

func (f *fileManagerImpl) GetIntConfigValueOrAskUser(key string, cmd *cobra.Command) (int, bool) {
	flag, ok := f.GetIntFlagValue(key, cmd)
	if ok {
		return flag, true
	}

	val, ok := f.GetIntConfigValue(key)
	if ok {
		return val, true
	}

	displayKey := strings.ReplaceAll(strings.ToUpper(key), "-", " ")
	return ToInt(f.ReadFromUser(displayKey))
}

func (f *fileManagerImpl) GetBoolConfigValueOrAskUser(key string, cmd *cobra.Command) (bool, bool) {
	flag, ok := f.GetBoolFlagValue(key, cmd)
	if ok {
		return flag, true
	}

	val, ok := f.GetBoolConfigValue(key)
	if ok {
		return val, true
	}

	displayKey := strings.ReplaceAll(strings.ToUpper(key), "-", " ")
	boolVal := f.ReadFromUser(displayKey)
	userVal, err := strconv.ParseBool(boolVal)
	if err != nil {
		return false, false
	}

	return userVal, true
}

func (f *fileManagerImpl) GetFlagOrConfigValue(key string, cmd *cobra.Command) (string, bool) {
	val, ok := f.GetFlagValue(key, cmd)
	if ok {
		return val, true
	}

	return f.GetConfigValue(key)
}

func (f *fileManagerImpl) GetIntFlagOrConfigValue(key string, cmd *cobra.Command) (int, bool) {
	val, ok := f.GetFlagOrConfigValue(key, cmd)
	if !ok {
		return 0, false
	}

	return ToInt(val)
}

func (f *fileManagerImpl) GetBoolFlagOrConfigValue(key string, cmd *cobra.Command) (bool, bool) {
	val, ok := f.GetFlagOrConfigValue(key, cmd)
	if !ok {
		return false, false
	}

	return ToBool(val)
}

func (f *fileManagerImpl) GetFlagValue(key string, cmd *cobra.Command) (string, bool) {
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

func (f *fileManagerImpl) GetIntFlagValue(key string, cmd *cobra.Command) (int, bool) {
	val, ok := f.GetFlagValue(key, cmd)
	if !ok {
		return 0, false
	}

	return ToInt(val)
}

func (f *fileManagerImpl) GetBoolFlagValue(key string, cmd *cobra.Command) (bool, bool) {
	val, ok := f.GetFlagValue(key, cmd)
	if !ok {
		return false, false
	}

	return ToBool(val)
}

func (f *fileManagerImpl) GetSensitiveConfigValueOrAskUser(key string, cmd *cobra.Command) (string, bool) {
	val, ok := f.GetFlagOrConfigValue(key, cmd)
	if ok {
		return val, true
	}

	val = f.ReadSensitiveFromUser(key)
	return val, val != ""
}

func (f *fileManagerImpl) GetFlagOrConfigValueOrAskUser(key string, cmd *cobra.Command) (string, bool) {
	val, ok := f.GetFlagOrConfigValue(key, cmd)
	if ok {
		return val, true
	}

	val = f.ReadFromUser(key)
	return val, val != ""
}

func (f *fileManagerImpl) GetSensitiveFlagOrConfigValueOrAskUser(key string, cmd *cobra.Command) (string, bool) {
	val, ok := f.GetFlagOrConfigValue(key, cmd)
	if ok {
		return val, true
	}

	val = f.ReadSensitiveFromUser(key)
	return val, val != ""
}

func (f *fileManagerImpl) GetIntFlagOrConfigValueOrAskUser(key string, cmd *cobra.Command) (int, bool) {
	val, ok := f.GetFlagOrConfigValue(key, cmd)
	if !ok {
		val = f.ReadFromUser(key)
	}

	return ToInt(val)
}

func (f *fileManagerImpl) GetBoolFlagOrConfigValueOrAskUser(key string, cmd *cobra.Command) (bool, bool) {
	val, ok := f.GetFlagOrConfigValue(key, cmd)
	if !ok {
		val = f.ReadFromUser(key)
	}

	return ToBool(val)
}
