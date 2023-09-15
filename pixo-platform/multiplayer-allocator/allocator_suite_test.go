package multiplayer_allocator

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestMultiplayerAllocator(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Multiplayer Allocator Suite")
}
