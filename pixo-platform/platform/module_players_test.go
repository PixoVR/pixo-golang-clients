package platform_test

import (
	"context"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("ModulePlayers", func() {

	It("can get module players", func() {
		modulePlayers, err := tokenClient.GetModulePlayers(context.Background())
		Expect(err).NotTo(HaveOccurred())
		Expect(modulePlayers).NotTo(BeNil())
		Expect(len(modulePlayers)).To(BeNumerically(">", 0))
		Expect(modulePlayers[0].ID).NotTo(BeZero())
		Expect(modulePlayers[0].Name).NotTo(BeEmpty())
	})

})
