package cmd_test

import (
	"fmt"
	"github.com/PixoVR/pixo-golang-server-utilities/pixo-platform/k8s/agones"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"io"
	"math/rand"
)

var _ = Describe("Deploy", func() {

	It("can deploy a server version", func() {
		majorVersion := rand.Intn(100)
		minorVersion := rand.Intn(100)
		patchVersion := rand.Intn(100)
		semanticVersion := fmt.Sprintf("%d.%d.%d", majorVersion, minorVersion, patchVersion)

		rootCmd, output := GetRootCmd()
		rootCmd.SetArgs([]string{
			"mp",
			"serverVersions",
			"deploy",
			"--module-id",
			"1",
			"--server-version",
			semanticVersion,
			"--image",
			agones.SimpleGameServerImage,
		})
		err := rootCmd.Execute()
		Expect(err).NotTo(HaveOccurred())

		out, err := io.ReadAll(output)
		Expect(err).NotTo(HaveOccurred())
		Expect(string(out)).To(ContainSubstring(fmt.Sprintf("created multiplayer server version: %s", semanticVersion)))
	})

	It("can check if a server version exists", func() {
		rootCmd, output := GetRootCmd()
		rootCmd.SetArgs([]string{
			"mp",
			"serverVersions",
			"deploy",
			"--pre-check",
			"--module-id",
			"1",
			"--server-version",
			"1.00.00",
		})
		err := rootCmd.Execute()
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("already exists"))

		rootCmd.SetArgs([]string{
			"mp",
			"serverVersions",
			"deploy",
			"--pre-check",
			"--module-id",
			"1",
			"--server-version",
			"99.99.99",
		})
		err = rootCmd.Execute()
		Expect(err).NotTo(HaveOccurred())

		out, err := io.ReadAll(output)
		Expect(err).NotTo(HaveOccurred())
		Expect(string(out)).To(ContainSubstring("does not exist"))
	})

})
