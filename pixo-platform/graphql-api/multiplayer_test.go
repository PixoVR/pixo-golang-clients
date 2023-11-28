package graphql_api_test

import (
	"fmt"
	. "github.com/PixoVR/pixo-golang-clients/pixo-platform/graphql-api"
	"github.com/PixoVR/pixo-golang-server-utilities/pixo-platform/k8s/agones"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"math/rand"
)

var _ = Describe("GraphQL API", func() {

	var (
		gqlClient *GraphQLAPIClient
	)

	BeforeEach(func() {
		gqlClient = NewClient("", "")
		Expect(gqlClient).NotTo(BeNil())
		Expect(gqlClient.IsAuthenticated()).To(BeTrue())
	})

	It("should be able to get the multiplayer server configs with a secret key", func() {
		mpServerConfigs, err := gqlClient.GetMultiplayerServerConfigs(MultiplayerServerConfigParams{
			OrgID:    1,
			ModuleID: 1,
		})
		Expect(err).NotTo(HaveOccurred())
		Expect(mpServerConfigs).NotTo(BeEmpty())
	})

	It("should be able to get the multiplayer server versions with a secret key", func() {
		mpServerVersions, err := gqlClient.GetMultiplayerServerVersions(MultiplayerServerVersionQueryParams{
			ModuleID:        1,
			SemanticVersion: "1.00.00",
		})
		Expect(err).NotTo(HaveOccurred())
		Expect(mpServerVersions).NotTo(BeEmpty())
	})

	It("should be able to create a multiplayer server version with a secret key", func() {
		randVersion := fmt.Sprintf("1.%d.%d", rand.Intn(100), rand.Intn(100))
		err := gqlClient.CreateMultiplayerServerVersion(1, agones.SimpleGameServerImage, randVersion)
		Expect(err).NotTo(HaveOccurred())
	})

})
