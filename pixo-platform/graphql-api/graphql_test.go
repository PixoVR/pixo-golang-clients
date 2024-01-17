package graphql_api_test

import (
	"context"
	"fmt"
	. "github.com/PixoVR/pixo-golang-clients/pixo-platform/graphql-api"
	platform "github.com/PixoVR/pixo-golang-clients/pixo-platform/primary-api"
	"github.com/PixoVR/pixo-golang-server-utilities/pixo-platform/k8s/agones"
	"github.com/go-faker/faker/v4"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"math/rand"
	"os"
)

var _ = Describe("GraphQL API", func() {

	var (
		secretKeyClient *GraphQLAPIClient
		tokenClient     *GraphQLAPIClient
		lifecycle       = "local"
		ctx             context.Context
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

		ctx = context.Background()
	})

	It("can login", func() {
		client := NewClient("", lifecycle, "na")
		err := client.Login(os.Getenv("PIXO_USERNAME"), os.Getenv("PIXO_PASSWORD"))
		Expect(err).NotTo(HaveOccurred())
		Expect(client.IsAuthenticated()).To(BeTrue())
	})

	It("can create a service account with a secret key and login with it", func() {
		user := platform.User{
			FirstName: faker.FirstName(),
			LastName:  faker.LastName(),
			Username:  faker.Username(),
			Password:  faker.Password(),
			OrgID:     1,
		}
		serviceAccount, err := tokenClient.CreateUser(ctx, user)
		Expect(err).NotTo(HaveOccurred())
		Expect(serviceAccount).NotTo(BeNil())
		Expect(serviceAccount.ID).NotTo(BeZero())
	})

	It("can create get and update a session, and then create an event with a secret key", func() {
		session, err := tokenClient.CreateSession(ctx, 1, "127.0.0.1", "test")
		Expect(err).NotTo(HaveOccurred())
		Expect(session).NotTo(BeNil())
		Expect(session.ID).NotTo(BeZero())

		retrievedSession, err := tokenClient.GetSession(ctx, session.ID)
		Expect(err).NotTo(HaveOccurred())
		Expect(retrievedSession).NotTo(BeNil())
		Expect(retrievedSession.ID).To(Equal(session.ID))
		Expect(retrievedSession.UserID).NotTo(BeZero())

		updatedSession, err := tokenClient.UpdateSession(ctx, session.ID, "TERMINATED", true)
		Expect(err).NotTo(HaveOccurred())
		Expect(updatedSession).NotTo(BeNil())
		Expect(updatedSession.ID).To(Equal(session.ID))

		event, err := tokenClient.CreateEvent(ctx, session.ID, "test", "test", "{}")
		Expect(err).NotTo(HaveOccurred())
		Expect(event).NotTo(BeNil())
	})

	It("can get the multiplayer server configs with a secret key", func() {
		mpServerConfigs, err := secretKeyClient.GetMultiplayerServerConfigs(ctx, MultiplayerServerConfigParams{
			OrgID:    1,
			ModuleID: 1,
		})
		Expect(err).NotTo(HaveOccurred())
		Expect(mpServerConfigs).NotTo(BeEmpty())
	})

	It("can get the multiplayer server versions with a secret key", func() {
		mpServerVersions, err := secretKeyClient.GetMultiplayerServerVersions(ctx, MultiplayerServerVersionQueryParams{
			ModuleID:        1,
			SemanticVersion: "1.00.00",
		})
		Expect(err).NotTo(HaveOccurred())
		Expect(mpServerVersions).NotTo(BeNil())
	})

	It("can create a multiplayer server version with a secret key", func() {
		randVersion := fmt.Sprintf("1.%d.%d", rand.Intn(100), rand.Intn(100))
		err := secretKeyClient.CreateMultiplayerServerVersion(ctx, 1, agones.SimpleGameServerImage, randVersion)
		Expect(err).NotTo(HaveOccurred())
	})

})
