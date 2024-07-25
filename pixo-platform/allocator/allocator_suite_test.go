package allocator

import (
	"github.com/PixoVR/pixo-golang-clients/pixo-platform/env"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestMultiplayerAllocator(t *testing.T) {
	RegisterFailHandler(Fail)
	env.SourceProjectEnv()
	RunSpecs(t, "Multiplayer Allocator Suite")
}
