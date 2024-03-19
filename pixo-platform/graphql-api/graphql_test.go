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
	"os"
)

var _ = Describe("GraphQL API", func() {

	var (
		ctx           context.Context
		moduleID      = 43
		orgID         = 20
		serverVersion = "1.03.02"
		randVersion   string
	)

	BeforeEach(func() {
		ctx = context.Background()
		randVersion = fmt.Sprintf("%d.%d.%d", rand.Intn(100), rand.Intn(100), rand.Intn(100))
	})

	It("can create a client, login, and interact with the api", func() {
		config := urlfinder.ClientConfig{Lifecycle: lifecycle, APIKey: pixoAPIKey}
		client := NewClient(config)
		Expect(client).NotTo(BeNil())

		Expect(client.Login(pixoUsername, pixoPassword)).To(Succeed())
		Expect(client.IsAuthenticated()).To(BeTrue())
		Expect(client.GetToken()).NotTo(BeEmpty())
		session, err := client.CreateSession(ctx, moduleID, "127.0.0.1", "test")
		Expect(err).NotTo(HaveOccurred())
		Expect(session).NotTo(BeNil())
		Expect(session.ID).NotTo(BeZero())
	})

	It("can create get and update a session, and then create an event with a secret key", func() {
		session, err := tokenClient.CreateSession(ctx, moduleID, "127.0.0.1", "test")
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
			ModuleID:      moduleID,
			OrgID:         orgID,
			ServerVersion: serverVersion,
		})
		Expect(err).NotTo(HaveOccurred())
		Expect(mpServerConfigs).NotTo(BeEmpty())
		Expect(mpServerConfigs[0].ServerVersions).NotTo(BeEmpty())
	})

	It("can get the multiplayer server versions", func() {
		mpServerVersions, err := tokenClient.GetMultiplayerServerVersions(ctx, &MultiplayerServerVersionQueryParams{
			ModuleID:        moduleID,
			SemanticVersion: serverVersion,
		})
		Expect(err).NotTo(HaveOccurred())
		Expect(mpServerVersions).NotTo(BeNil())
		Expect(mpServerVersions).NotTo(BeEmpty())
	})

	It("can return an error if no image or file are given", func() {
		_, err := tokenClient.CreateMultiplayerServerVersion(ctx, MultiplayerServerVersion{})
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("image or file path must be provided"))
	})

	It("can create and get a multiplayer server version", func() {
		input := MultiplayerServerVersion{
			ModuleID:        moduleID,
			SemanticVersion: randVersion,
			ImageRegistry:   agones.SimpleGameServerImage,
			Engine:          "unreal",
		}

		serverVersion, err := tokenClient.CreateMultiplayerServerVersion(ctx, input)

		Expect(err).NotTo(HaveOccurred())
		Expect(serverVersion).NotTo(BeNil())

		mpServerVersion, err := tokenClient.GetMultiplayerServerVersion(ctx, serverVersion.ID)

		Expect(err).NotTo(HaveOccurred())
		Expect(mpServerVersion).NotTo(BeNil())
		Expect(mpServerVersion.ID).To(Equal(serverVersion.ID))
		Expect(mpServerVersion.ModuleID).To(Equal(moduleID))
		Expect(mpServerVersion.SemanticVersion).To(Equal(randVersion))
		Expect(mpServerVersion.ImageRegistry).To(Equal(agones.SimpleGameServerImage))
		Expect(mpServerVersion.Engine).To(Equal(input.Engine))
		Expect(mpServerVersion.Status).To(Equal("enabled"))
	})

	It("can upload a gameserver build", func() {
		localFilePath := "./test.zip"
		file, err := os.Create(localFilePath)
		Expect(err).NotTo(HaveOccurred())
		n, err := file.WriteString("test")
		defer func() {
			_ = file.Close()
			_ = os.Remove(localFilePath)
		}()
		Expect(err).NotTo(HaveOccurred())
		Expect(n).To(Equal(4))
		serverVersionInput := MultiplayerServerVersion{
			ModuleID:        moduleID,
			SemanticVersion: randVersion,
			Engine:          "unreal",
			LocalFilePath:   localFilePath,
		}

		serverVersion, err := tokenClient.CreateMultiplayerServerVersion(ctx, serverVersionInput)

		Expect(err).NotTo(HaveOccurred())
		Expect(serverVersion).NotTo(BeNil())
		Expect(serverVersion.ID).NotTo(BeZero())
		Expect(serverVersion.FileLink).NotTo(BeEmpty())
		Expect(serverVersion.FileLink).To(ContainSubstring("X-Goog-Algorithm"))
		Expect(serverVersion.ImageRegistry).To(BeEmpty())
	})

})
