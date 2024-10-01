package platform_test

import (
	"context"
	"fmt"
	"github.com/PixoVR/pixo-golang-clients/pixo-platform/platform"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"math/rand"
	"os"
)

var _ = Describe("Assets", func() {

	var (
		ctx context.Context
	)

	BeforeEach(func() {
		ctx = context.Background()
	})

	It("can return an error if passed a nil pointer when creating an asset", func() {
		Expect(tokenClient.CreateAsset(ctx, nil)).To(MatchError("asset is nil"))
	})

	It("can return an error if passed a nil pointer when creating an asset version", func() {
		Expect(tokenClient.CreateAssetVersion(ctx, nil)).To(MatchError("asset version is nil"))
	})

	It("can return an error if request fails", func() {
		Expect(tokenClient.CreateAsset(ctx, &platform.Asset{})).To(MatchError("module id or external id is required"))
	})

	It("can return an error if request fails", func() {
		Expect(tokenClient.CreateAssetVersion(ctx, &platform.AssetVersion{})).To(MatchError("asset id is required"))
	})

	Context("managing assets", func() {

		var (
			assetName string
			asset     *platform.Asset
		)

		BeforeEach(func() {
			assetName = fmt.Sprintf("logo-%d", rand.Intn(1000000))
			asset = &platform.Asset{
				ModuleID: moduleID,
				Name:     assetName,
				Type:     "image",
			}
			Expect(tokenClient.CreateAsset(ctx, asset)).To(Succeed())
			Expect(asset).NotTo(BeNil())
			Expect(asset.ID).NotTo(BeZero())
		})

		It("can get an asset", func() {
			asset, err := tokenClient.GetAsset(ctx, asset.ID)
			Expect(err).NotTo(HaveOccurred())
			Expect(asset).NotTo(BeNil())
			Expect(asset.ID).To(Equal(asset.ID))
			Expect(asset.ModuleID).To(Equal(moduleID))
			Expect(asset.Name).To(Equal(assetName))
		})

		It("can get assets", func() {
			asset, err := tokenClient.GetAssets(ctx, platform.AssetParams{ModuleID: moduleID})
			Expect(err).NotTo(HaveOccurred())
			Expect(len(asset)).To(BeNumerically(">", 0))
		})

		Context("with a file", func() {

			var localFilePath = "./test.png"

			BeforeEach(func() {
				Expect(os.WriteFile(localFilePath, []byte("test"), 0644)).NotTo(HaveOccurred())
			})

			AfterEach(func() {
				Expect(os.Remove(localFilePath)).To(Succeed())
			})

			It("can create an asset version", func() {
				assetVersion := platform.AssetVersion{
					AssetID:       asset.ID,
					LanguageCode:  "en",
					LocalFilePath: localFilePath,
				}

				Expect(tokenClient.CreateAssetVersion(ctx, &assetVersion)).To(Succeed())

				Expect(assetVersion).NotTo(BeNil())
				Expect(assetVersion.ID).NotTo(BeZero())
				Expect(assetVersion.FileLink).NotTo(BeEmpty())
				Expect(assetVersion.FileLink).To(ContainSubstring("X-Goog-Algorithm"))
			})

			It("can update an asset versions status", func() {
				assetVersion := platform.AssetVersion{
					AssetID:       asset.ID,
					LanguageCode:  "en",
					LocalFilePath: localFilePath,
				}
				Expect(tokenClient.CreateAssetVersion(ctx, &assetVersion)).To(Succeed())

				assetVersion.Status = "release"
				Expect(tokenClient.UpdateAssetVersion(ctx, &assetVersion)).To(Succeed())
			})

		})

	})

})
