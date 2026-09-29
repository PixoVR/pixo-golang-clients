package platform_test

import (
	"context"
	"time"

	"github.com/go-faker/faker/v4"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/PixoVR/pixo-golang-clients/pixo-platform/platform"
)

var _ = Describe("OrgModules", func() {

	var (
		ctx context.Context
		org *platform.Org
	)

	BeforeEach(func() {
		ctx = context.Background()

		var err error
		org, err = tokenClient.CreateOrg(ctx, platform.Org{
			AffiliateID: 20,
			Name:        faker.Username(),
			Type:        "customer",
			Status:      "enabled",
		})
		Expect(err).NotTo(HaveOccurred())
		Expect(org.ID).NotTo(BeZero())

		expiresAt := time.Now().AddDate(0, 0, 3)
		orgModule, err := tokenClient.CreateOrgModule(ctx, platform.OrgModule{
			OrgID:     org.ID,
			ModuleID:  moduleID,
			ExpirDate: &expiresAt,
		})
		Expect(err).NotTo(HaveOccurred())
		Expect(orgModule.ID).NotTo(BeZero())
	})

	AfterEach(func() {
		Expect(tokenClient.DeleteOrgModule(ctx, org.ID, moduleID)).To(Succeed())
	})

	It("can return an error if the org id is missing", func() {
		orgModules, err := tokenClient.GetOrgModules(ctx, platform.OrgModuleParams{})
		Expect(err).To(HaveOccurred())
		Expect(orgModules).To(BeNil())
	})

	It("can return the modules an org has access to", func() {
		orgModules, err := tokenClient.GetOrgModules(ctx, platform.OrgModuleParams{OrgID: org.ID})

		Expect(err).NotTo(HaveOccurred())
		Expect(orgModules).NotTo(BeEmpty())

		orgModule := orgModules[0]
		Expect(orgModule.OrgID).To(Equal(org.ID))
		Expect(orgModule.Module).NotTo(BeNil())
		Expect(orgModule.Module.ID).To(Equal(moduleID))
		Expect(orgModule.Module.Abbreviation).NotTo(BeEmpty())
		Expect(orgModule.ExpiresAt).NotTo(BeEmpty())
	})
})
