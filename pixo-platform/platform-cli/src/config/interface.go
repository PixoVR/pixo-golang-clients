package config

import (
	"github.com/spf13/cobra"
)

type Manager interface {
	GetConfig() *Config
	SetConfig(config Config)
}

type ValueFinder interface {
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
}
