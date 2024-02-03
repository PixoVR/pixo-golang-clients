package config

import (
	"github.com/spf13/cobra"
	"io"
)

type Manager interface {
	SetReader(r io.Reader)
	Reader() io.Reader
	SetWriter(w io.Writer)
	Writer() io.Writer

	GetActiveEnv() Env
	SetActiveEnv(env Env)
	ConfigFile() string
	SetConfigFile(configFile string) error

	GetConfig() Config
	SetConfig(config Config)

	ReadFromUser(key string) string
	ReadFromUserOrReturn(key, defaultValue string) string
	ReadSensitiveFromUser(key string) string

	GetFlagValue(key string, cmd *cobra.Command) (string, bool)
	GetIntFlagValue(key string, cmd *cobra.Command) (int, bool)
	GetBoolFlagValue(key string, cmd *cobra.Command) (bool, bool)

	GetConfigValue(key string) (string, bool)
	GetIntConfigValue(key string) (int, bool)
	GetBoolConfigValue(key string) (bool, bool)

	GetFlagOrConfigValue(key string, cmd *cobra.Command) (string, bool)
	GetIntFlagOrConfigValue(key string, cmd *cobra.Command) (int, bool)
	GetBoolFlagOrConfigValue(key string, cmd *cobra.Command) (bool, bool)

	SetConfigValue(key, value string)
	SetIntConfigValue(key string, value int)
	SetBoolConfigValue(key string, value bool)

	UnsetConfigValue(key string)

	GetConfigValueOrAskUser(key string, cmd *cobra.Command) (string, bool)
	GetSensitiveConfigValueOrAskUser(key string, cmd *cobra.Command) (string, bool)
	GetIntConfigValueOrAskUser(key string, cmd *cobra.Command) (int, bool)
	GetBoolConfigValueOrAskUser(key string, cmd *cobra.Command) (bool, bool)

	GetFlagOrConfigValueOrAskUser(key string, cmd *cobra.Command) (string, bool)
	GetSensitiveFlagOrConfigValueOrAskUser(key string, cmd *cobra.Command) (string, bool)
	GetIntFlagOrConfigValueOrAskUser(key string, cmd *cobra.Command) (int, bool)
	GetBoolFlagOrConfigValueOrAskUser(key string, cmd *cobra.Command) (bool, bool)

	Lifecycle() string
	Region() string
}
