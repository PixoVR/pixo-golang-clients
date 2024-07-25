package config

import (
	"fmt"
	"os"
	"strings"
)

type Config struct {
	Region    string         `mapstructure:"region"`
	Lifecycle string         `mapstructure:"lifecycle"`
	Envs      map[string]Env `mapstructure:"envs"`
}

func (c *Config) ActiveEnv() Env {
	if c.Envs == nil {
		c.Envs = map[string]Env{}
	}

	if c.Region == "" {
		c.Region = "na"
	}

	if c.Lifecycle == "" {
		c.Lifecycle = "prod"
	}

	env, ok := c.Envs[envName(c.Region, c.Lifecycle)]
	if !ok {
		env = Env{
			Region:    c.Region,
			Lifecycle: c.Lifecycle,
			EnvMap:    map[string]string{},
		}
	}

	return env
}

func (c *Config) SetEnv(env Env) {
	if c.Envs == nil {
		c.Envs = map[string]Env{}
	}

	c.Envs[env.Name()] = env
}

func lookupConfigEnv(key string) (string, bool) {
	envVarName := formatEnvVarName(key)
	return os.LookupEnv(envVarName)
}

func formatEnvVarName(key string) string {
	return fmt.Sprintf("PIXO_%s", strings.ReplaceAll(strings.ToUpper(key), "-", "_"))
}
