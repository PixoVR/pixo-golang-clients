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

	//It("should be able to create a multiplayer server version with a secret key", func() {
	//multiplayerServerVersion := primary_api.MultiplayerServerVersion{
	//	ModuleID:         moduleID,
	//	Status:           "enabled",
	//	ImageRegistry:    image,
	//	Engine:           "unreal",
	//	Version:          semanticVersion,
	//	MinClientVersion: semanticVersion,
	//}

	//	image := "us-docker.pkg.dev/agones-images/examples/simple-game-server:0.14"
	//	res, err := gqlClient.CreateMultiplayerServerVersion(125, image, "1.00.00")
	//	Expect(err).NotTo(HaveOccurred())
	//	Expect(res.StatusCode()).To(Equal(http.StatusOK))
	//})

})
