package cmd_test

import (
	"bytes"
	"errors"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Module", func() {

	BeforeEach(func() {
		executor = NewTestExecutor()
	})

	AfterEach(func() {
		executor.Cleanup()
	})

	It("can create a module version", func() {
		input := bytes.NewBufferString("1\n")

		output, err := executor.RunCommandWithInput(
			input,
			"modules",
			"deploy",
		)

		Expect(err).NotTo(HaveOccurred())
		Expect(output).To(ContainSubstring("Enter MODULE ID:"))
		Expect(output).To(ContainSubstring("Enter SEMANTIC VERSION:"))
		Expect(output).To(ContainSubstring("Semantic version not provided"))
	})

	It("can return an error if module id is missing", func() {
		input := bytes.NewBufferString("0\n")

		output, err := executor.RunCommandWithInput(
			input,
			"modules",
			"deploy",
			"--module-id",
			"0",
		)

		Expect(err).NotTo(HaveOccurred())
		Expect(output).To(ContainSubstring("Enter MODULE ID:"))
		Expect(output).To(ContainSubstring("Module ID not provided"))
		Expect(executor.MockMatchmakingClient.NumCalledDialWebsocket).To(Equal(0))
		Expect(executor.MockMatchmakingClient.NumCalledWriteToWebsocketError).To(Equal(0))
		Expect(executor.MockMatchmakingClient.NumCalledReadFromWebsocket).To(Equal(0))
		Expect(executor.MockMatchmakingClient.NumCalledCloseWebsocket).To(Equal(0))
	})

	It("can return an error if semantic version is missing", func() {
		input := bytes.NewBufferString("\n")

		output, err := executor.RunCommandWithInput(
			input,
			"modules",
			"deploy",
			"--module-id",
			"1",
			"--package",
			"pixovr.com",
			"--platforms",
			"1",
			"--controls",
			"1",
			"--zip-file",
			"test.zip",
		)

		Expect(err).NotTo(HaveOccurred())
		Expect(output).To(ContainSubstring("Enter SEMANTIC VERSION:"))
		Expect(output).To(ContainSubstring("Semantic version not provided"))
	})

	It("can return an error if package name is missing", func() {
		input := bytes.NewBufferString("")

		output, err := executor.RunCommandWithInput(
			input,
			"modules",
			"deploy",
			"--module-id",
			"1",
			"--semantic-version",
			"1.0.0",
		)

		Expect(err).NotTo(HaveOccurred())
		Expect(output).To(ContainSubstring("Enter PACKAGE:"))
		Expect(output).To(ContainSubstring("Package name not provided"))
	})

	It("can return an error if zip file path is missing", func() {
		input := bytes.NewBufferString("")

		output, err := executor.RunCommandWithInput(
			input,
			"modules",
			"deploy",
			"--module-id",
			"1",
			"--semantic-version",
			"1.0.0",
			"--package",
			"test",
			"--zip-file",
			"",
		)

		Expect(err).NotTo(HaveOccurred())
		Expect(output).To(ContainSubstring("Enter ZIP FILE:"))
		Expect(output).To(ContainSubstring("Zip file not provided"))
	})

	It("can return an error if the platform options cant be found", func() {
		executor.MockPlatformClient.GetPlatformsError = errors.New("error")
		input := bytes.NewBufferString("")

		output, err := executor.RunCommandWithInput(
			input,
			"modules",
			"deploy",
			"--module-id",
			"1",
			"--semantic-version",
			"1.0.0",
			"--package",
			"test",
			"--zip-file",
			"test.zip",
		)

		Expect(err).NotTo(HaveOccurred())
		Expect(output).To(ContainSubstring("error"))
		Expect(executor.MockPlatformClient.NumCalledGetPlatforms).To(Equal(1))
	})

	It("can return an error if platforms are missing", func() {
		input := bytes.NewBufferString("")

		output, err := executor.RunCommandWithInput(
			input,
			"modules",
			"deploy",
			"--module-id",
			"1",
			"--semantic-version",
			"1.0.0",
			"--package",
			"test",
			"--zip-file",
			"test.zip",
		)

		Expect(err).NotTo(HaveOccurred())
		Expect(executor.MockPlatformClient.NumCalledGetPlatforms).To(Equal(1))
		Expect(output).To(ContainSubstring("Select PLATFORMS:"))
		Expect(output).To(ContainSubstring("Platforms not provided"))
	})

	It("can return an error if the control types options cant be found", func() {
		executor.MockPlatformClient.GetControlTypesError = errors.New("error")
		input := bytes.NewBufferString("1\n")

		output, err := executor.RunCommandWithInput(
			input,
			"modules",
			"deploy",
			"--module-id",
			"1",
			"--semantic-version",
			"1.0.0",
			"--package",
			"test",
			"--zip-file",
			"test.zip",
		)

		Expect(err).NotTo(HaveOccurred())
		Expect(executor.MockPlatformClient.NumCalledGetControlTypes).To(Equal(1))
		Expect(output).To(ContainSubstring("error"))
	})

	It("can return an error if control types are missing", func() {
		input := bytes.NewBufferString("1\n")

		output, err := executor.RunCommandWithInput(
			input,
			"modules",
			"deploy",
			"--module-id",
			"1",
			"--semantic-version",
			"1.0.0",
			"--package",
			"test",
			"--zip-file",
			"test.zip",
		)

		Expect(err).NotTo(HaveOccurred())
		Expect(executor.MockPlatformClient.NumCalledGetControlTypes).To(Equal(1))
		Expect(output).To(ContainSubstring("Select CONTROL TYPES:"))
		Expect(output).To(ContainSubstring("Control types not provided"))
	})

	It("can return an error if the api call fails", func() {
		executor.MockPlatformClient.CreateModuleVersionError = errors.New("error")
		input := bytes.NewBufferString("1\n1\n")

		output, err := executor.RunCommandWithInput(
			input,
			"modules",
			"deploy",
			"--module-id",
			"1",
			"--semantic-version",
			"1.0.0",
			"--package",
			"test",
			"--zip-file",
			"test.zip",
			"--platforms",
			"1",
			"--controls",
			"1",
		)

		Expect(err).NotTo(HaveOccurred())
		Expect(output).To(ContainSubstring("error"))
		Expect(executor.MockPlatformClient.NumCalledCreateModuleVersion).To(Equal(1))
	})

	It("can deploy a module version", func() {
		output, err := executor.RunCommand(
			"modules",
			"deploy",
			"--module-id",
			"1",
			"--semantic-version",
			"1.0.0",
			"--package",
			"test",
			"--zip-file",
			"test.zip",
			"--platforms",
			"1",
			"--controls",
			"1",
		)

		Expect(err).NotTo(HaveOccurred())
		Expect(output).To(ContainSubstring("Deployed version 1.0.0 for module 1"))
		Expect(executor.MockPlatformClient.NumCalledCreateModuleVersion).To(Equal(1))
	})

})
