package urlfinder_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestURLFinder(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "URL Finder Suite")
}
