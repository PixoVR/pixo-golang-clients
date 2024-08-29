package platform_test

import (
	"context"
	"fmt"
	"github.com/PixoVR/pixo-golang-clients/pixo-platform/allocator"
	"github.com/PixoVR/pixo-golang-clients/pixo-platform/platform"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/rs/zerolog/log"
	"math/rand"
	"os"
)

var _ = Describe("Multiplayer Resources", func() {

	var (
		ctx           context.Context
		serverVersion = "1.03.02"
		randVersion   string
		localFilePath = "./test.zip"
	)

	BeforeEach(func() {
		ctx = context.Background()
		randVersion = fmt.Sprintf("%d.%d.%d", rand.Intn(100), rand.Intn(100), rand.Intn(100))
	})

	It("can get the multiplayer server configs", func() {
		mpServerConfigs, err := tokenClient.GetMultiplayerServerConfigs(ctx, &platform.MultiplayerServerConfigParams{
			ModuleID:      moduleID,
			OrgID:         orgID,
			ServerVersion: serverVersion,
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

	It("can get the multiplayer server versions", func() {
		mpServerVersions, err := tokenClient.GetMultiplayerServerVersions(ctx, &platform.MultiplayerServerVersionQueryParams{
			ModuleID:        moduleID,
			SemanticVersion: serverVersion,
		})
		Expect(err).NotTo(HaveOccurred())
		Expect(mpServerVersions).NotTo(BeNil())
		Expect(mpServerVersions).NotTo(BeEmpty())
	})

	It("can create and get a multiplayer server version", func() {
		input := platform.MultiplayerServerVersion{
			ModuleID:        moduleID,
			SemanticVersion: randVersion,
			ImageRegistry:   allocator.SimpleGameServerImage,
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
		Expect(mpServerVersion.ImageRegistry).To(Equal(allocator.SimpleGameServerImage))
		Expect(mpServerVersion.Engine).To(Equal(input.Engine))
		Expect(mpServerVersion.Status).To(Equal("enabled"))
	})

	It("can upload a gameserver build", func() {
		file, err := os.Create(localFilePath)
		Expect(err).NotTo(HaveOccurred())
		n, err := file.WriteString("test")
		defer func() {
			_ = file.Close()
			_ = os.Remove(localFilePath)
		}()
		Expect(err).NotTo(HaveOccurred())
		Expect(n).To(Equal(4))
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
