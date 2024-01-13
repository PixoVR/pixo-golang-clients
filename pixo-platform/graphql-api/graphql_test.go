package graphql_api_test

import (
	"fmt"
	. "github.com/PixoVR/pixo-golang-clients/pixo-platform/graphql-api"
	"github.com/PixoVR/pixo-golang-server-utilities/pixo-platform/k8s/agones"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"math/rand"
	"os"
)

var _ = Describe("GraphQL API", func() {

	var (
		secretKeyClient *GraphQLAPIClient
		tokenClient     *GraphQLAPIClient
		lifecycle       = "dev"
	)

	BeforeEach(func() {
		secretKeyClient = NewClient("", lifecycle, "")
		Expect(secretKeyClient).NotTo(BeNil())
		Expect(secretKeyClient.IsAuthenticated()).To(BeTrue())

		var err error
		tokenClient, err = NewClientWithBasicAuth(os.Getenv("PIXO_USERNAME"), os.Getenv("PIXO_PASSWORD"), lifecycle, "")
		Expect(err).NotTo(HaveOccurred())
		Expect(tokenClient).NotTo(BeNil())
		Expect(tokenClient.IsAuthenticated()).To(BeTrue())
		Expect(tokenClient.GetToken()).NotTo(BeEmpty())
	})

	It("should be able to login", func() {
		client := NewClient("", lifecycle, "na")
		err := client.Login(os.Getenv("PIXO_USERNAME"), os.Getenv("PIXO_PASSWORD"))
		Expect(err).NotTo(HaveOccurred())
		Expect(client.IsAuthenticated()).To(BeTrue())
	})

	It("should be able to get create and get a session with a secret key", func() {
		session, err := tokenClient.CreateSession(1, "127.0.0.1", "test")
		Expect(err).NotTo(HaveOccurred())
		Expect(session).NotTo(BeNil())
		Expect(session.ID).NotTo(BeZero())

		retrievedSession, err := tokenClient.GetSession(session.ID)
		Expect(err).NotTo(HaveOccurred())
		Expect(retrievedSession).NotTo(BeNil())
		Expect(retrievedSession.ID).To(Equal(session.ID))
		Expect(retrievedSession.UserID).NotTo(BeZero())
	})

	It("should be able to get the multiplayer server configs with a secret key", func() {
		mpServerConfigs, err := secretKeyClient.GetMultiplayerServerConfigs(MultiplayerServerConfigParams{
			OrgID:    1,
			ModuleID: 1,
		})
		Expect(err).NotTo(HaveOccurred())
		Expect(mpServerConfigs).NotTo(BeEmpty())
	})

	It("should be able to get the multiplayer server versions with a secret key", func() {
		mpServerVersions, err := secretKeyClient.GetMultiplayerServerVersions(MultiplayerServerVersionQueryParams{
			ModuleID:        1,
			SemanticVersion: "1.00.00",
		})
		Expect(err).NotTo(HaveOccurred())
		Expect(mpServerVersions).NotTo(BeNil())
	})

	It("should be able to create a multiplayer server version with a secret key", func() {
		randVersion := fmt.Sprintf("1.%d.%d", rand.Intn(100), rand.Intn(100))
		err := secretKeyClient.CreateMultiplayerServerVersion(1, agones.SimpleGameServerImage, randVersion)
		Expect(err).NotTo(HaveOccurred())
	})

})
