package graphql_api_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"os"

	. "github.com/PixoVR/pixo-golang-clients/pixo-platform/graphql-api"
)

var _ = Describe("GraphQL API", func() {

	var (
		gqlClient *GraphQLAPIClient
	)

	BeforeEach(func() {
		gqlClient = NewClient(os.Getenv("SECRET_KEY"), "")
		Expect(gqlClient).NotTo(BeNil())
		Expect(gqlClient.IsAuthenticated()).To(BeTrue())
	})

	It("should be able to get the multiplayer server versions with a secret key", func() {
		mpServerVersions, err := gqlClient.GetMultiplayerServerVersions()
		Expect(err).NotTo(HaveOccurred())
		Expect(mpServerVersions).NotTo(BeNil())
		Expect(len(mpServerVersions)).To(BeNumerically(">", 0))
	})

	It("should be able to create a multiplayer server version with a secret key", func() {
		image := "us-docker.pkg.dev/agones-images/examples/simple-game-server:0.14"
		err := gqlClient.DeployMultiplayerServerVersion(1, image, "1.00.00")
		Expect(err).NotTo(HaveOccurred())
	})

	It("should be able to get a matchmaking profile by org id and module id", func() {
		Expect(1).To(Equal(2))
	})

})
