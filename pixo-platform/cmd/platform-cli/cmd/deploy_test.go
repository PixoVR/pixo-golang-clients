package cmd_test

import (
	"fmt"
	"github.com/PixoVR/pixo-golang-server-utilities/pixo-platform/k8s/agones"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"math/rand"
)

var _ = Describe("Deploy", Ordered, func() {

	var (
		semanticVersion string
	)

	BeforeAll(func() {
		semanticVersion = fmt.Sprintf("%d.%d.%d", rand.Intn(100), rand.Intn(100), rand.Intn(100))
	})

	It("can deploy a server version", func() {
		output, err := RunCommand(
			"mp",
			"serverVersions",
			"deploy",
			"--module-id",
			"1",
			"--server-version",
			semanticVersion,
			"--image",
			agones.SimpleGameServerImage,
		)
		Expect(err).NotTo(HaveOccurred())
		Expect(output).To(ContainSubstring(fmt.Sprintf("created multiplayer server version: %s", semanticVersion)))
	})

	It("can check if a server version exists", func() {
		_, err := RunCommand(
			"mp",
			"serverVersions",
			"deploy",
			"--pre-check",
			"--module-id",
			"1",
			"--server-version",
			semanticVersion,
		)
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("already exists"))

		output, err := RunCommand(
			"mp",
			"serverVersions",
			"deploy",
			"--pre-check",
			"--module-id",
			"1",
			"--server-version",
			"99.99.99",
		)
		Expect(err).NotTo(HaveOccurred())
		Expect(output).To(ContainSubstring("does not exist"))
	})

})
