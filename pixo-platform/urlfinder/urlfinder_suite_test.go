package urlfinder_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestUrlfinder(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Urlfinder Suite")
}
