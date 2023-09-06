package abstract_client_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestAbstractClient(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "AbstractClient Suite")
}
