package platform_test

import (
	"context"
	"time"

	platform "github.com/PixoVR/pixo-golang-clients/pixo-platform/platform"
	faker "github.com/go-faker/faker/v4"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Reporting", func() {
	var (
		ctx      context.Context
		org      *platform.Org
		moduleID = 43
	)

	BeforeEach(func() {
		ctx = context.Background()

		var err error
		org, err = tokenClient.CreateOrg(ctx, platform.Org{
			AffiliateID: 20,
			Name:        faker.Username(),
			Type:        "distributor",
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

	It("can return expiring modules", func() {
		expiringModules, err := tokenClient.GetExpiringModules(ctx)

		Expect(err).NotTo(HaveOccurred())
		Expect(expiringModules).NotTo(BeEmpty())

		var found bool
		for _, expiringModule := range expiringModules {
			if expiringModule.OrgID == org.ID {
				found = true
				Expect(expiringModule.OrgName).To(Equal(org.Name))
				Expect(expiringModule.ExpiresAt).To(BeTemporally("~", time.Now().AddDate(0, 0, 3), time.Hour))
			}
		}
		Expect(found).To(BeTrue())
	})
})
