package platform_test

import (
	"context"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Reporting", func() {
	var (
		ctx = context.Background()
	)

	It("can return expiring modules", func() {
		expiringModules, err := tokenClient.GetExpiringModules(ctx)
		Expect(err).NotTo(HaveOccurred())
		Expect(expiringModules).NotTo(BeNil())
	})
})
