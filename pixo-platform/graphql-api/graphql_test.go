package graphql_api_test

import (
	"context"
	"fmt"
	. "github.com/PixoVR/pixo-golang-clients/pixo-platform/graphql-api"
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

	It("can get the multiplayer server configs", func() {
		mpServerConfigs, err := tokenClient.GetMultiplayerServerConfigs(ctx, MultiplayerServerConfigParams{
			OrgID:    1,
			ModuleID: 1,
		})
		Expect(err).NotTo(HaveOccurred())
		Expect(mpServerConfigs).NotTo(BeEmpty())
	})

	It("can get the multiplayer server versions", func() {
		mpServerVersions, err := tokenClient.GetMultiplayerServerVersions(ctx, MultiplayerServerVersionQueryParams{
			ModuleID:        1,
			SemanticVersion: "1.00.00",
		})
		Expect(err).NotTo(HaveOccurred())
		Expect(mpServerVersions).NotTo(BeNil())
	})

	It("can create a multiplayer server version", func() {
		randVersion := fmt.Sprintf("1.%d.%d", rand.Intn(100), rand.Intn(100))
		err := tokenClient.CreateMultiplayerServerVersion(ctx, 1, agones.SimpleGameServerImage, randVersion)
		Expect(err).NotTo(HaveOccurred())
	})

})
