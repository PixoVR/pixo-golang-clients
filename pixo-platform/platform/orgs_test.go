package platform_test

import (
	"context"
	platform "github.com/PixoVR/pixo-golang-clients/pixo-platform/platform"
	"github.com/go-faker/faker/v4"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Orgs API", func() {

	var (
		ctx      context.Context
		orgInput platform.Org
		testOrg  *platform.Org
	)

	BeforeEach(func() {
		ctx = context.Background()
		orgInput = platform.Org{
			Name:       faker.Username(),
			Type:       "distributor",
			OpenAccess: false,
		}
		var err error

		testOrg, err = tokenClient.CreateOrg(ctx, orgInput)

		Expect(err).NotTo(HaveOccurred())
		Expect(testOrg).NotTo(BeNil())
		Expect(testOrg.ID).NotTo(BeZero())
	})

	//AfterEach(func() {
	//	Expect(tokenClient.DeleteOrg(ctx, testOrg.ID)).To(Succeed())
	//	deletedOrg, err := tokenClient.GetOrgByID(ctx, testOrg.ID)
	//	Expect(err).To(HaveOccurred())
	//	Expect(deletedOrg).To(BeNil())
	//})

	It("can get an org by ID", func() {
		retrievedOrg, err := tokenClient.GetOrg(ctx, testOrg.ID)

		Expect(err).NotTo(HaveOccurred())
		Expect(retrievedOrg).NotTo(BeNil())
		Expect(retrievedOrg.ID).To(Equal(testOrg.ID))
		Expect(retrievedOrg.HubLogoLink).NotTo(BeEmpty())
	})

	//It("can get all orgs", func() {
	//	orgs, err := tokenClient.GetOrgs(ctx)
	//
	//	Expect(err).NotTo(HaveOccurred())
	//	Expect(orgs).NotTo(BeNil())
	//	Expect(len(orgs)).To(BeNumerically(">", 0))
	//})

})
