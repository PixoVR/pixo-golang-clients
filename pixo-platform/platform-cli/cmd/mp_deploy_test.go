package cmd_test

import (
	"bytes"
	"errors"
	"fmt"
	. "github.com/PixoVR/pixo-golang-clients/pixo-platform/allocator"
	"github.com/kyokomi/emoji"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"math/rand"
	"os"
)

var _ = Describe("Server Deploy", func() {

	var (
		semanticVersion = fmt.Sprintf("%d.%d.%d", rand.Intn(100), rand.Intn(100), rand.Intn(100))
	)

	BeforeEach(func() {
		executor = NewTestExecutor()
	})

	AfterEach(func() {
		executor.Cleanup()
	})

	It("should return an error if the modules get call fails", func() {
		executor.MockPlatformClient.GetModulesError = errors.New("failed to get modules")

		output, err := executor.RunCommand(
			"mp",
			"servers",
			"deploy",
			"--pre-check",
		)

		Expect(err).To(HaveOccurred())
		Expect(err).To(MatchError("failed to get modules"))
		Expect(output).To(BeEmpty())
	})

	It("should ask for the module id and server version if it is not provided", func() {
		input := bytes.NewBufferString("1: TST - test\n")

		output, err := executor.RunCommandWithInput(
			input,
			"mp",
			"servers",
			"deploy",
			"--pre-check",
		)

		Expect(err).To(HaveOccurred())
		Expect(err).To(MatchError("SERVER VERSION not provided"))
		Expect(output).To(ContainSubstring("MODULE ID"))
		Expect(output).To(ContainSubstring("SERVER VERSION"))
	})

	It("can tell if a server version exists", func() {
		output, err := executor.RunCommand(
			"mp",
			"servers",
			"deploy",
			"--pre-check",
			"--module-id",
			"1: TST - test",
			"--server-version",
			"1.00.00",
		)
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("server version 1.00.00 already exists"))
		Expect(output).To(ContainSubstring(emoji.Sprint("\b:exclamation: server version 1.00.00 already exists\n")))
	})

	It("can tell if a server version does not exist", func() {
		executor.MockPlatformClient.GetMultiplayerServerVersionsEmpty = true
		output := executor.RunCommandAndExpectSuccess(
			"mp",
			"servers",
			"deploy",
			"--pre-check",
			"--module-id",
			"1: TST - test",
			"--server-version",
			"99.99.99",
		)
		Expect(output).To(Equal(emoji.Sprint("\b:heavy_check_mark: Server version does not exist yet: 99.99.99\n")))
	})

	It("can ask the user for the docker image to use if not provided", func() {
		input := bytes.NewBufferString("test\n")

		output := executor.RunCommandWithInputAndExpectSuccess(
			input,
			"mp",
			"servers",
			"deploy",
			"--module-id",
			"1: TST - test",
			"--server-version",
			semanticVersion,
		)

		Expect(output).To(ContainSubstring("Enter DOCKER IMAGE:"))
		Expect(output).To(ContainSubstring(fmt.Sprintf("Deployed version: %s", semanticVersion)))
	})

	It("can deploy a server version", func() {
		output := executor.RunCommandAndExpectSuccess(
			"mp",
			"servers",
			"deploy",
			"--module-id",
			"1: TST - test",
			"--server-version",
			semanticVersion,
			"--image",
			SimpleGameServerImage,
		)
		Expect(output).To(ContainSubstring(fmt.Sprintf("Deployed version: %s", semanticVersion)))
	})

	It("can upload a server version with a zip file", func() {
		localFilePath := "./test.zip"
		file, err := os.Create(localFilePath)
		Expect(err).NotTo(HaveOccurred())
		_, err = file.WriteString("test")
		Expect(err).NotTo(HaveOccurred())
		defer func() {
			_ = file.Close()
			_ = os.Remove(localFilePath)
		}()

		output, err := executor.RunCommand(
			"mp",
			"servers",
			"deploy",
			"--module-id",
			"1: TST - test",
			"--server-version",
			semanticVersion,
			"--zip-file",
			localFilePath,
		)

		Expect(err).NotTo(HaveOccurred())
		Expect(output).To(ContainSubstring(fmt.Sprintf("Deployed version: %s", semanticVersion)))
		Expect(executor.MockPlatformClient.NumCalledCreateMultiplayerServerVersion).To(Equal(1))
	})

})
