package primary_api_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestPrimaryApi(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "PrimaryApi Suite")
}
