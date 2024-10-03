package platform_test

import (
	"context"
	"fmt"
	"github.com/PixoVR/pixo-golang-clients/pixo-platform/allocator"
	"github.com/PixoVR/pixo-golang-clients/pixo-platform/platform"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"math/rand"
	"os"
)

var _ = Describe("Multiplayer Resources", func() {

	var (
		ctx             context.Context
		semanticVersion = "1.03.02"
		randVersion     string
		localFilePath   = "./test.zip"
	)

	BeforeEach(func() {
		ctx = context.Background()
		randVersion = fmt.Sprintf("%d.%d.%d", rand.Intn(100), rand.Intn(100), rand.Intn(100))
	})

	It("can get the multiplayer server configs", func() {
		mpServerConfigs, err := tokenClient.GetMultiplayerServerConfigs(ctx, &platform.MultiplayerServerConfigParams{
			ModuleID:      moduleID,
			OrgID:         orgID,
			ServerVersion: semanticVersion,
		})
		Expect(err).NotTo(HaveOccurred())
		Expect(mpServerConfigs).NotTo(BeEmpty())
		Expect(mpServerConfigs[0].ServerVersions).NotTo(BeEmpty())
	})

	It("can return an error if no image or file are given", func() {
		_, err := tokenClient.CreateMultiplayerServerVersion(ctx, platform.MultiplayerServerVersion{})
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("image or file path must be provided"))
	})

	It("can get the multiplayer server versions with a config", func() {
		mpServerVersions, err := tokenClient.GetMultiplayerServerVersionsWithConfig(ctx, &platform.MultiplayerServerVersionParams{
			ModuleID:        moduleID,
			SemanticVersion: semanticVersion,
		})
		Expect(err).NotTo(HaveOccurred())
		Expect(len(mpServerVersions)).To(BeNumerically(">", 0))
		for _, mpServerVersion := range mpServerVersions {
			Expect(mpServerVersion.ID).NotTo(BeZero())
			Expect(mpServerVersion.ModuleID).To(Equal(moduleID))
			Expect(mpServerVersion.SemanticVersion).To(Equal(semanticVersion))
			Expect(mpServerVersion.ImageRegistry).NotTo(BeEmpty())
		}
	})

	It("can get multiplayer server versions", func() {
		mpServerVersions, err := tokenClient.GetMultiplayerServerVersions(ctx, nil)
		Expect(err).NotTo(HaveOccurred())
		Expect(len(mpServerVersions)).To(BeNumerically(">", 0))
		for _, mpServerVersion := range mpServerVersions {
			Expect(mpServerVersion.ID).NotTo(BeZero())
			Expect(mpServerVersion.ModuleID).NotTo(BeZero())
			Expect(mpServerVersion.SemanticVersion).NotTo(BeEmpty())
		}
	})

	Context("managing server versions", func() {

		var (
			serverVersion *platform.MultiplayerServerVersion
			input         platform.MultiplayerServerVersion
		)

		BeforeEach(func() {
			input = platform.MultiplayerServerVersion{
				ModuleID:        moduleID,
				SemanticVersion: randVersion,
				ImageRegistry:   allocator.SimpleGameServerImage,
				Engine:          "unreal",
			}

			var err error
			serverVersion, err = tokenClient.CreateMultiplayerServerVersion(ctx, input)

			Expect(err).NotTo(HaveOccurred())
			Expect(serverVersion).NotTo(BeNil())
		})

		It("can get a multiplayer server version by id", func() {
			retrievedServerVersion, err := tokenClient.GetMultiplayerServerVersion(ctx, serverVersion.ID)

			Expect(err).NotTo(HaveOccurred())
			Expect(retrievedServerVersion).NotTo(BeNil())
			Expect(retrievedServerVersion.ID).To(Equal(serverVersion.ID))
			Expect(retrievedServerVersion.ModuleID).To(Equal(moduleID))
			Expect(retrievedServerVersion.SemanticVersion).To(Equal(randVersion))
			Expect(retrievedServerVersion.ImageRegistry).To(Equal(input.ImageRegistry))
			Expect(retrievedServerVersion.Engine).To(Equal(input.Engine))
			Expect(retrievedServerVersion.Status).To(Equal("enabled"))
		})

		It("can update a multiplayer server version", func() {
			input.ImageRegistry = "gcr.io/pixo-bootstrap/multiplayer/gameservers/simple-server:latest"
			input.Status = "disabled"

			updatedServerVersion, err := tokenClient.UpdateMultiplayerServerVersion(ctx, input)

			Expect(err).NotTo(HaveOccurred())
			Expect(updatedServerVersion).NotTo(BeNil())
			Expect(updatedServerVersion.ID).To(Equal(serverVersion.ID))
			Expect(updatedServerVersion.ModuleID).To(Equal(moduleID))
			Expect(updatedServerVersion.SemanticVersion).To(Equal(randVersion))
			Expect(updatedServerVersion.ImageRegistry).To(Equal(input.ImageRegistry))
			Expect(updatedServerVersion.Engine).To(Equal(input.Engine))
			Expect(updatedServerVersion.Status).To(Equal(input.Status))
		})

	})

	It("can upload a gameserver build", func() {
		cleanup := NewTestFile(localFilePath)
		defer cleanup()
		serverVersionInput := platform.MultiplayerServerVersion{
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

func NewTestFile(filePath string) func() {
	Expect(os.WriteFile(filePath, []byte("test"), 0644)).To(Succeed())
	return func() {
		Expect(os.Remove(filePath)).To(Succeed())
	}
}
