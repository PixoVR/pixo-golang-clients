package platform_test

import (
	"context"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Roles", func() {

	It("can retrieve roles", func() {
		roles, err := tokenClient.GetRoles(context.Background())
		Expect(err).NotTo(HaveOccurred())
		Expect(roles).NotTo(BeNil())
		Expect(len(roles)).To(BeNumerically(">", 0))
	})

})
