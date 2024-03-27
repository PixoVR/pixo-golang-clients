package cmd_test

import (
	"bytes"
	"fmt"
	"github.com/PixoVR/pixo-golang-server-utilities/pixo-platform/k8s/agones"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"math/rand"
	"os"
)

var _ = Describe("Deploy", Ordered, func() {

	var (
		executor        *TestExecutor
		semanticVersion string
	)

	BeforeAll(func() {
		semanticVersion = fmt.Sprintf("%d.%d.%d", rand.Intn(100), rand.Intn(100), rand.Intn(100))
		executor = NewTestExecutor()
		Expect(executor).NotTo(BeNil())
		Expect(semanticVersion).NotTo(BeEmpty())
	})

	AfterAll(func() {
		executor.Cleanup()
	})

	It("can deploy a server version", func() {
		output, err := executor.RunCommand(
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
		Expect(output).To(ContainSubstring(fmt.Sprintf("Deployed version: %s", semanticVersion)))
	})

	It("can tell if a server version exists", func() {
		_, err := executor.RunCommand(
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
	})

	It("can tell if a server version does not exist", func() {
		executor.MockPlatformClient.GetMultiplayerServerVersionsEmpty = true
		output, err := executor.RunCommand(
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

	It("can ask the user for the docker image to use if not provided", func() {
		input := bytes.NewBufferString("test\n")

		output, err := executor.RunCommandWithInput(
			input,
			"mp",
			"serverVersions",
			"deploy",
			"--module-id",
			"1",
			"--server-version",
			semanticVersion,
		)

		Expect(err).NotTo(HaveOccurred())
		Expect(output).To(ContainSubstring("Enter DOCKER IMAGE:"))
		Expect(output).To(ContainSubstring(fmt.Sprintf("Deployed version: %s", semanticVersion)))
	})

	It("can upload a server version with a zip file", func() {
		localFilePath := "./test.zip"
		file, err := os.Create(localFilePath)
		Expect(err).NotTo(HaveOccurred())
		defer func() {
			_ = file.Close()
			_ = os.Remove(localFilePath)
		}()
		_, err = file.WriteString("test")
		Expect(err).NotTo(HaveOccurred())

		output, err := executor.RunCommand(
			"mp",
			"serverVersions",
			"deploy",
			"--module-id",
			"1",
			"--server-version",
			semanticVersion,
			"--zip-file",
			localFilePath,
		)

		Expect(err).NotTo(HaveOccurred())
		Expect(output).To(ContainSubstring(fmt.Sprintf("Deployed version: %s", semanticVersion)))
		Expect(executor.MockPlatformClient.CalledCreateMultiplayerServerVersion).To(BeTrue())
	})

})
