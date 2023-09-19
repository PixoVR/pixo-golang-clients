package graphql_api_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestGraphQLAPISuite(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "GraphQL API Suite")
}
