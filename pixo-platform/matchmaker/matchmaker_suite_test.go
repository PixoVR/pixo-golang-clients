package matchmaker_test

import (
	"github.com/PixoVR/pixo-golang-clients/pixo-platform/env"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestMultiplayer(t *testing.T) {
	RegisterFailHandler(Fail)
	env.SourceProjectEnv()
	RunSpecs(t, "Multiplayer Suite")
}
