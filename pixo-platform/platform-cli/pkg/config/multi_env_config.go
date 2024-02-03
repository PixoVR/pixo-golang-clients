package config

import (
	"gopkg.in/yaml.v3"
	"os"
)

type Config struct {
	Region    string         `mapstructure:"region"`
	Lifecycle string         `mapstructure:"lifecycle"`
	Envs      map[string]Env `mapstructure:"envs"`
}

func (f *Config) ActiveEnv() Env {
	if f.Envs == nil {
		f.Envs = map[string]Env{}
	}

	if f.Region == "" {
		f.Region = "na"
	}

	if f.Lifecycle == "" {
		f.Lifecycle = "prod"
	}

	env, ok := f.Envs[envName(f.Region, f.Lifecycle)]
	if !ok {
		env = Env{
			Region:    f.Region,
			Lifecycle: f.Lifecycle,
		}
	}

	return env
}

func (f *Config) SetEnv(env Env) {
	if f.Envs == nil {
		f.Envs = map[string]Env{}
	}

	f.Envs[env.Name()] = env
}

func WriteToFile(path string, configFile Config) error {
	configFileBytes, err := yaml.Marshal(configFile)
	if err != nil {
		return err
	}

	if err = os.WriteFile(path, configFileBytes, 0644); err != nil {
		return err
	}

	return nil
}
