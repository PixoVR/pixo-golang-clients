package config_test

import (
	"github.com/PixoVR/pixo-golang-clients/pixo-platform/platform-cli/src/config"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestConfig(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Config Suite")
}

var (
	sampleConfig = config.Config{
		Lifecycle: "prod",
		Region:    "na",
		Envs: map[string]config.Env{
			"na-prod": {
				Region:    "na",
				Lifecycle: "prod",
				EnvMap: map[string]string{
					"api-key":  "na-prod-api-key",
					"token":    "na-prod-token",
					"username": "na-prod-username",
					"password": "na-prod-password",
				},
			},
			"na-stage": {
				Region:    "na",
				Lifecycle: "stage",
				EnvMap: map[string]string{
					"api-key":  "na-stage-api-key",
					"token":    "na-stage-token",
					"username": "na-stage-username",
					"password": "na-stage-password",
				},
			},
		},
	}
)
