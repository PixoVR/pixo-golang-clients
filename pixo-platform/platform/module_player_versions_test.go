package platform_test

import (
	"context"

	"github.com/PixoVR/pixo-golang-clients/pixo-platform/platform"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("ModulePlayerVersions", func() {

	var ctx = context.Background()

	It("can get all module player versions without params", func() {
		versions, err := tokenClient.GetModulePlayerVersions(ctx, nil)
		Expect(err).NotTo(HaveOccurred())
		Expect(versions).NotTo(BeNil())
	})

	It("can get module player versions filtered by status", func() {
		params := &platform.ModulePlayerVersionParams{
			Status: []string{"enabled"},
		}
		versions, err := tokenClient.GetModulePlayerVersions(ctx, params)
		Expect(err).NotTo(HaveOccurred())
		Expect(versions).NotTo(BeNil())
		for _, v := range versions {
			Expect(v.Status).To(Equal("enabled"))
		}
	})

})
