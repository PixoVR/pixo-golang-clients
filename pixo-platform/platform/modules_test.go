package platform_test

import (
	"context"
	"fmt"
	"math/rand"

	. "github.com/PixoVR/pixo-golang-clients/pixo-platform/platform"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
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
		cleanup := NewTestFile(localFilePath)
		defer cleanup()
		_, err := tokenClient.CreateModuleVersion(ctx, ModuleVersion{LocalFilePath: localFilePath})
		Expect(err).To(HaveOccurred())
	})

	It("can get modules", func() {
		modules, err := tokenClient.GetModules(ctx)
		Expect(err).NotTo(HaveOccurred())
		Expect(modules).NotTo(BeEmpty())
		Expect(modules[0].ID).NotTo(BeZero())
		Expect(modules[0].Abbreviation).NotTo(BeEmpty())
		Expect(modules[0].Status).NotTo(BeEmpty())
	})

	It("can get only the public modules", func() {
		isPublic := true

		modules, err := tokenClient.GetModules(ctx, ModuleParams{IsPublic: &isPublic})

		Expect(err).NotTo(HaveOccurred())
		Expect(modules).NotTo(BeEmpty())
		for _, module := range modules {
			Expect(module.IsPublic).To(BeTrue())
		}
	})

	It("can get only the modules of one distributor", func() {
		module, err := tokenClient.GetModule(ctx, moduleID)
		Expect(err).NotTo(HaveOccurred())
		Expect(module.DistributorID).NotTo(BeZero())

		modules, err := tokenClient.GetModules(ctx, ModuleParams{DistributorID: &module.DistributorID})

		Expect(err).NotTo(HaveOccurred())
		Expect(modules).NotTo(BeEmpty())
		for _, m := range modules {
			Expect(m.DistributorID).To(Equal(module.DistributorID))
		}
	})

	It("can get only the modules with the given status", func() {
		modules, err := tokenClient.GetModules(ctx, ModuleParams{Statuses: []string{"enabled"}})

		Expect(err).NotTo(HaveOccurred())
		Expect(modules).NotTo(BeEmpty())
		for _, module := range modules {
			Expect(module.Status).To(Equal("enabled"))
		}
	})

	It("can return an error if getting a module without an id", func() {
		module, err := tokenClient.GetModule(ctx, 0)
		Expect(err).To(HaveOccurred())
		Expect(module).To(BeNil())
	})

	It("can get a module with the fields the headset needs", func() {
		module, err := tokenClient.GetModule(ctx, moduleID)

		Expect(err).NotTo(HaveOccurred())
		Expect(module.ID).To(Equal(moduleID))
		Expect(module.Abbreviation).NotTo(BeEmpty())
		Expect(module.Status).NotTo(BeEmpty())
		Expect(module.DistributorID).NotTo(BeZero())
		Expect(module.Distributor).NotTo(BeNil())
		Expect(module.Distributor.Name).NotTo(BeEmpty())
	})

	It("can get a module with its versions", func() {
		module, err := tokenClient.GetModule(ctx, moduleID)

		Expect(err).NotTo(HaveOccurred())
		Expect(module.Versions).NotTo(BeEmpty())
		version := module.Versions[0]
		Expect(version.ID).NotTo(BeZero())
		Expect(version.ModuleID).To(Equal(moduleID))
		Expect(version.SemanticVersion).NotTo(BeEmpty())
		Expect(version.LifecycleID).NotTo(BeZero())
		Expect(version.CreatedAt).NotTo(BeZero())
	})

	It("can get module versions with the fields the headset needs", func() {
		versions, err := tokenClient.GetModuleVersions(ctx, &ModuleVersionParams{ModuleID: &moduleID})

		Expect(err).NotTo(HaveOccurred())
		Expect(versions).NotTo(BeEmpty())
		Expect(versions[0].SemanticVersion).NotTo(BeEmpty())
		Expect(versions[0].Lifecycle).NotTo(BeNil())
		Expect(versions[0].Module.Abbreviation).NotTo(BeEmpty())
		Expect(versions[0].CreatedAt).NotTo(BeZero())
	})

	It("can create a module version", func() {
		cleanup := NewTestFile(localFilePath)
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
		Expect(moduleVersion.LifecycleID).To(Equal(1))
		Expect(moduleVersion.Lifecycle).NotTo(BeNil())
		Expect(moduleVersion.Lifecycle.Name).To(Equal(LifecycleDevelopment))
		Expect(moduleVersion.FileLink).NotTo(BeEmpty())
	})

})
