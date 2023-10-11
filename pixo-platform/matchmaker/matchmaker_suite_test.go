package matchmaker_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestMultiplayer(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Multiplayer Suite")
}
