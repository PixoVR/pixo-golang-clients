package config

import (
	"fmt"
)

type Env struct {
	EnvMap    map[string]string `mapstructure:"envmap"`
	Region    string            `mapstructure:"region"`
	Lifecycle string            `mapstructure:"lifecycle"`
}

func (e *Env) Name() string {
	return envName(e.Region, e.Lifecycle)
}

func envName(region, lifecycle string) string {
	return fmt.Sprintf("%s-%s", region, lifecycle)
}

func (e *Env) Get(key string) (string, bool) {
	if e.EnvMap == nil {
		e.EnvMap = map[string]string{}
	}

	val, ok := e.EnvMap[key]
	return val, ok
}

func (e *Env) Set(key, value string) {
	if e.EnvMap == nil {
		e.EnvMap = map[string]string{}
	}

	e.EnvMap[key] = value
}

func (e *Env) Unset(key string) {
	if e.EnvMap == nil {
		e.EnvMap = map[string]string{}
	}

	delete(e.EnvMap, key)
}
