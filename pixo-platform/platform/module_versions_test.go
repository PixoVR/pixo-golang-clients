package platform_test

import (
	"github.com/PixoVR/pixo-golang-clients/pixo-platform/platform"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"context"
)

var _ = Describe("ModuleVersions", func() {
	var ctx = context.Background()

	It("can get all module versions without params", func() {
		platforms, err := tokenClient.GetModuleVersions(ctx, nil)
		Expect(err).NotTo(HaveOccurred())
		Expect(platforms).NotTo(BeEmpty())
	})

	It("can get all module versions with params", func() {
		params := &platform.ModuleVersionParams{
			Lifecycles:         []string{platform.LifecycleReleased},
			ModuleID:           &moduleID,
			PlatformShortNames: []string{"android"},
			SortField:          nil,
			SortOrder:          nil,
		}
		platforms, err := tokenClient.GetModuleVersions(ctx, params)
		Expect(err).NotTo(HaveOccurred())
		Expect(platforms).NotTo(BeNil())
	})

})
