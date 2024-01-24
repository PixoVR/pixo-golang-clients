package graphql_api_test

import (
	"github.com/PixoVR/pixo-golang-clients/pixo-platform/urlfinder"
	"os"
	"testing"

	. "github.com/PixoVR/pixo-golang-clients/pixo-platform/graphql-api"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestGraphQLAPISuite(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "GraphQL API Suite")
}

var (
	secretKeyClient *GraphQLAPIClient
	tokenClient     *GraphQLAPIClient
	lifecycle       = "local"
	pixoUsername    = os.Getenv("PIXO_USERNAME")
	pixoPassword    = os.Getenv("PIXO_PASSWORD")
)

var _ = BeforeSuite(func() {
	config := urlfinder.ClientConfig{Lifecycle: lifecycle}
	secretKeyClient = NewClient(config)
	Expect(secretKeyClient).NotTo(BeNil())
	Expect(secretKeyClient.IsAuthenticated()).To(BeTrue())

	var err error
	tokenClient, err = NewClientWithBasicAuth(pixoUsername, pixoPassword, config)
	Expect(err).NotTo(HaveOccurred())
	Expect(tokenClient).NotTo(BeNil())
	Expect(tokenClient.IsAuthenticated()).To(BeTrue())
	Expect(tokenClient.GetToken()).NotTo(BeEmpty())
})
