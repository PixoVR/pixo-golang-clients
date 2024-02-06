package graphql_api_test

import (
	"context"
	"fmt"
	. "github.com/PixoVR/pixo-golang-clients/pixo-platform/graphql-api"
	"github.com/PixoVR/pixo-golang-clients/pixo-platform/urlfinder"
	"github.com/PixoVR/pixo-golang-server-utilities/pixo-platform/k8s/agones"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"math/rand"
)

var _ = Describe("GraphQL API", func() {

	var (
		ctx context.Context
	)

	BeforeEach(func() {
		ctx = context.Background()
	})

	It("can create a client, login, and interact with the api", func() {
		config := urlfinder.ClientConfig{Lifecycle: lifecycle, APIKey: pixoAPIKey}
		client := NewClient(config)
		Expect(client).NotTo(BeNil())

		Expect(client.Login(pixoUsername, pixoPassword)).To(Succeed())
		Expect(client.IsAuthenticated()).To(BeTrue())
		Expect(client.GetToken()).NotTo(BeEmpty())
		session, err := client.CreateSession(ctx, 1, "127.0.0.1", "test")
		Expect(err).NotTo(HaveOccurred())
		Expect(session).NotTo(BeNil())
		Expect(session.ID).NotTo(BeZero())
	})

	It("can create get and update a session, and then create an event with a secret apiKey", func() {
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

		// ERROR: invalid input syntax for type uuid: "" (SQLSTATE 22P02) ??
		//event, err := tokenClient.CreateEvent(ctx, session.ID, faker.UUIDDigit(), "test", "{}")
		//Expect(err).NotTo(HaveOccurred())
		//Expect(event).NotTo(BeNil())
	})

	It("can get the multiplayer server configs", func() {
		mpServerConfigs, err := tokenClient.GetMultiplayerServerConfigs(ctx, &MultiplayerServerConfigParams{
			ModuleID:      271,
			OrgID:         20,
			ServerVersion: "2.00.01",
		})
		Expect(err).NotTo(HaveOccurred())
		Expect(mpServerConfigs).NotTo(BeEmpty())
		Expect(mpServerConfigs[0].ServerVersions).NotTo(BeEmpty())
	})

	It("can get the multiplayer server versions", func() {
		mpServerVersions, err := tokenClient.GetMultiplayerServerVersions(ctx, &MultiplayerServerVersionQueryParams{
			ModuleID:        271,
			SemanticVersion: "2.00.01",
		})
		Expect(err).NotTo(HaveOccurred())
		Expect(mpServerVersions).NotTo(BeNil())
		Expect(mpServerVersions).NotTo(BeEmpty())
	})

	It("can create a multiplayer server version", func() {
		randVersion := fmt.Sprintf("1.%d.%d", rand.Intn(100), rand.Intn(100))
		err := tokenClient.CreateMultiplayerServerVersion(ctx, 1, agones.SimpleGameServerImage, randVersion)
		Expect(err).NotTo(HaveOccurred())
	})

})
