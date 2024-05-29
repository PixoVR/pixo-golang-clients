package graphql_api_test

import (
	"github.com/PixoVR/pixo-golang-clients/pixo-platform/urlfinder"
	config2 "github.com/PixoVR/pixo-golang-server-utilities/pixo-platform/config"
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
	apiKeyClient *GraphQLAPIClient
	tokenClient  *GraphQLAPIClient
	lifecycle    = config2.GetEnvOrReturn("PIXO_LIFECYCLE", "stage")
	pixoUsername = os.Getenv("PIXO_USERNAME")
	pixoPassword = os.Getenv("PIXO_PASSWORD")
	pixoAPIKey   = os.Getenv("PIXO_API_KEY")
)

var _ = BeforeSuite(func() {
	config := urlfinder.ClientConfig{Lifecycle: lifecycle, APIKey: pixoAPIKey}
	apiKeyClient = NewClient(config)
	Expect(apiKeyClient).NotTo(BeNil())
	Expect(apiKeyClient.IsAuthenticated()).To(BeTrue())

	var err error
	tokenClient, err = NewClientWithBasicAuth(pixoUsername, pixoPassword, config)
	Expect(err).NotTo(HaveOccurred())
	Expect(tokenClient).NotTo(BeNil())
	Expect(tokenClient.IsAuthenticated()).To(BeTrue())
	Expect(tokenClient.GetToken()).NotTo(BeEmpty())
})
