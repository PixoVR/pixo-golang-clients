package graphql_api_test

import (
	"context"
	"fmt"
	. "github.com/PixoVR/pixo-golang-clients/pixo-platform/graphql-api"
	"github.com/PixoVR/pixo-golang-clients/pixo-platform/urlfinder"
	"github.com/PixoVR/pixo-golang-server-utilities/pixo-platform/k8s/agones"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/rs/zerolog/log"
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
		localFilePath = "./test.zip"
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

	It("can create get and update a session", func() {
		session, err := tokenClient.CreateSession(ctx, moduleID, "127.0.0.1", "test")
		Expect(err).NotTo(HaveOccurred())
		Expect(session).NotTo(BeNil())
		Expect(session.ID).NotTo(BeZero())

		retrievedSession, err := tokenClient.GetSession(ctx, session.ID)
		Expect(err).NotTo(HaveOccurred())
		Expect(retrievedSession).NotTo(BeNil())
		Expect(retrievedSession.ID).To(Equal(session.ID))
		Expect(retrievedSession.UserID).NotTo(BeZero())

		input := Session{
			ID:        session.ID,
			Status:    "TERMINATED",
			Completed: true,
			RawScore:  0.5,
			MaxScore:  1.0,
		}

		updatedSession, err := tokenClient.UpdateSession(ctx, input)
		Expect(err).NotTo(HaveOccurred())
		Expect(updatedSession).NotTo(BeNil())
		Expect(updatedSession.ID).To(Equal(session.ID))
		Expect(updatedSession.UserID).To(Equal(retrievedSession.UserID))
		Expect(updatedSession.User.OrgID).To(Equal(retrievedSession.User.OrgID))
		Expect(updatedSession.ModuleID).To(Equal(retrievedSession.ModuleID))
		Expect(updatedSession.RawScore).To(Equal(input.RawScore))
		Expect(updatedSession.MaxScore).To(Equal(input.MaxScore))
		Expect(updatedSession.ScaledScore).To(Equal(input.RawScore / input.MaxScore))
		Expect(updatedSession.CompletedAt).NotTo(BeNil())
		Expect(updatedSession.Duration).NotTo(BeNil())
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

	It("can return an error if a required field is missing", func() {
		cleanup := makeTestFile(localFilePath)
		defer cleanup()
		_, err := tokenClient.CreateModuleVersion(ctx, ModuleVersion{LocalFilePath: localFilePath})
		Expect(err).To(HaveOccurred())
	})

	It("can get all platforms", func() {
		platforms, err := tokenClient.GetPlatforms(ctx)
		Expect(err).NotTo(HaveOccurred())
		Expect(platforms).NotTo(BeEmpty())
	})

	It("can get all control types", func() {
		controlTypes, err := tokenClient.GetControlTypes(ctx)
		Expect(err).NotTo(HaveOccurred())
		Expect(controlTypes).NotTo(BeEmpty())
	})

	It("can create a module version", func() {
		cleanup := makeTestFile(localFilePath)
		defer cleanup()
		input := ModuleVersion{
			ModuleID:        moduleID,
			LocalFilePath:   "./test.zip",
			SemanticVersion: "1.0.0",
			Package:         "test",
			PlatformIds:     []int{1},
			ControlIds:      []int{1},
		}

		moduleVersion, err := tokenClient.CreateModuleVersion(ctx, input)

		Expect(err).NotTo(HaveOccurred())
		Expect(moduleVersion).NotTo(BeNil())
		Expect(moduleVersion.ID).NotTo(BeZero())
		Expect(moduleVersion.ModuleID).To(Equal(moduleID))
		Expect(moduleVersion.SemanticVersion).To(Equal("1.0.0"))
		Expect(moduleVersion.Package).To(Equal("test"))
		Expect(moduleVersion.Status).To(Equal("disabled"))
		Expect(moduleVersion.FileLink).NotTo(BeEmpty())
	})

})

func makeTestFile(filePath string) func() {
	file, err := os.Create(filePath)
	if err != nil {
		log.Panic().Err(err).Msg("failed to create file")
	}

	if _, err = file.WriteString("test"); err != nil {
		log.Panic().Err(err).Msg("failed to write to file")
	}

	return func() {
		_ = file.Close()
		_ = os.Remove(filePath)
	}
}
