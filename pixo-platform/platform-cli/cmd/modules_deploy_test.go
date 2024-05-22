package cmd_test

import (
	"bytes"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Module", func() {

	var (
		executor *TestExecutor
	)

	BeforeEach(func() {
		executor = NewTestExecutor()
	})

	AfterEach(func() {
		executor.Cleanup()
	})

	It("can create a module version", func() {
		input := bytes.NewReader([]byte("1\n"))

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
		input := bytes.NewReader([]byte("0\n"))

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
		Expect(executor.MockMatchmakingClient.NumCalledDialMatchmaker).To(Equal(0))
		Expect(executor.MockMatchmakingClient.NumCalledWriteToWebsocket).To(Equal(0))
		Expect(executor.MockMatchmakingClient.NumCalledReadFromWebsocket).To(Equal(0))
		Expect(executor.MockMatchmakingClient.NumCalledCloseWebsocket).To(Equal(0))
	})

	It("can return an error if server version is missing", func() {
		input := bytes.NewReader([]byte("\n"))

		output, err := executor.RunCommandWithInput(
			input,
			"modules",
			"deploy",
			"--module-id",
			"1",
			"--semantic-version",
			"",
		)

		Expect(err).NotTo(HaveOccurred())
		Expect(output).To(ContainSubstring("Semantic version not provided"))
	})

	It("can return an error if package name is missing", func() {
		input := bytes.NewReader([]byte(""))

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
		input := bytes.NewReader([]byte(""))

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
		executor.MockPlatformClient.GetPlatformsError = true
		input := bytes.NewReader([]byte(""))

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
		Expect(executor.MockPlatformClient.CalledGetPlatforms).To(BeTrue())
	})

	It("can return an error if platforms are missing", func() {
		input := bytes.NewReader([]byte(""))

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
		Expect(executor.MockPlatformClient.CalledGetPlatforms).To(BeTrue())
		Expect(output).To(ContainSubstring("Select PLATFORMS:"))
		Expect(output).To(ContainSubstring("Platforms not provided"))
	})

	It("can return an error if the control types options cant be found", func() {
		executor.MockPlatformClient.GetControlTypesError = true
		input := bytes.NewReader([]byte("1\n"))

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
		Expect(executor.MockPlatformClient.CalledGetControlTypes).To(BeTrue())
		Expect(output).To(ContainSubstring("error"))
	})

	It("can return an error if control types are missing", func() {
		input := bytes.NewReader([]byte("1\n"))

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
		Expect(executor.MockPlatformClient.CalledGetControlTypes).To(BeTrue())
		Expect(output).To(ContainSubstring("Select CONTROL TYPES:"))
		Expect(output).To(ContainSubstring("Control types not provided"))
	})

	It("can return an error if the api call fails", func() {
		executor.MockPlatformClient.CreateModuleVersionError = true
		input := bytes.NewReader([]byte("1\n1\n"))

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
		Expect(executor.MockPlatformClient.CalledCreateModuleVersion).To(BeTrue())
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
		Expect(executor.MockPlatformClient.CalledCreateModuleVersion).To(BeTrue())
	})

})
