package basic_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestBasicForms(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Basic Forms Suite")
}
