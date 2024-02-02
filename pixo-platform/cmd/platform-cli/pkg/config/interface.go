package config

import "io"

type Manager interface {
	SetReader(r io.Reader)
	Reader() io.Reader
	SetWriter(w io.Writer)
	Writer() io.Writer

	GetActiveEnv() Env
	SetActiveEnv(env Env)
	ConfigFile() string
	ReadConfigFile(configFile string) error

	GetConfig() Config
	SetConfig(config Config)
	GetConfigValue(key string) (string, bool)
	SetConfigValue(key, value string)
	UnsetConfigValue(key, value string)

	GetConfigValueOrAskUser(key, shortFlag string) string

	Lifecycle() string
	Region() string
}
