package platform_test

import (
	"context"
	"fmt"
	. "github.com/PixoVR/pixo-golang-clients/pixo-platform/platform"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"math/rand"
)

var _ = Describe("Modules", func() {

	var (
		ctx           context.Context
		randVersion   string
		localFilePath = "./test.zip"
	)

	BeforeEach(func() {
		ctx = context.Background()
		randVersion = fmt.Sprintf("%d.%d.%d", rand.Intn(100), rand.Intn(100), rand.Intn(100))
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

	It("can return an error if a required field is missing when creating a module version", func() {
		cleanup := makeTestFile(localFilePath)
		defer cleanup()
		_, err := tokenClient.CreateModuleVersion(ctx, ModuleVersion{LocalFilePath: localFilePath})
		Expect(err).To(HaveOccurred())
	})

	It("can create a module version", func() {
		cleanup := makeTestFile(localFilePath)
		defer cleanup()
		input := ModuleVersion{
			ModuleID:        moduleID,
			LocalFilePath:   "./test.zip",
			SemanticVersion: randVersion,
			Package:         "test",
			PlatformIds:     []int{1},
			ControlIds:      []int{1},
		}

		moduleVersion, err := tokenClient.CreateModuleVersion(ctx, input)

		Expect(err).NotTo(HaveOccurred())
		Expect(moduleVersion).NotTo(BeNil())
		Expect(moduleVersion.ID).NotTo(BeZero())
		Expect(moduleVersion.ModuleID).To(Equal(moduleID))
		Expect(moduleVersion.SemanticVersion).To(Equal(randVersion))
		Expect(moduleVersion.Package).To(Equal(input.Package))
		Expect(moduleVersion.Status).To(Equal("disabled"))
		Expect(moduleVersion.FileLink).NotTo(BeEmpty())
	})

})
