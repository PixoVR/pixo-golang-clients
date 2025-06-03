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

	It("can return an error if the semantic version is not provided", func() {
		input := bytes.NewBufferString("TST\n")

		output, err := executor.RunCommandWithInput(
			input,
			"modules",
			"deploy",
		)

		Expect(err).To(MatchError("SEMANTIC VERSION not provided"))
		Expect(output).To(ContainSubstring("MODULE"))
		Expect(output).To(ContainSubstring("SEMANTIC VERSION"))
	})

	It("can return an error if module is missing", func() {
		input := bytes.NewBufferString("\n")

		output, err := executor.RunCommandWithInput(
			input,
			"modules",
			"deploy",
			"--package",
			"package",
			"--semantic-version",
			"1.0.0",
			"--zip-file",
			"test.zip",
			"--platforms",
			"android",
			"--controls",
			"keyboard/mouse",
		)

		Expect(err).To(MatchError("MODULE not provided"))
		Expect(output).To(ContainSubstring("MODULE"))
	})

	It("can return an error if semantic version is missing", func() {
		input := bytes.NewBufferString("")

		output, err := executor.RunCommandWithInput(
			input,
			"modules",
			"deploy",
			"--module",
			"TST",
			"--package",
			"pixovr.com",
			"--platforms",
			"android",
			"--controls",
			"keyboard/mouse",
			"--zip-file",
			"test.zip",
		)

		Expect(err).To(MatchError("SEMANTIC VERSION not provided"))
		Expect(output).To(ContainSubstring("SEMANTIC VERSION"))
	})

	It("can return an error if package name is missing", func() {
		input := bytes.NewBufferString("")

		output, err := executor.RunCommandWithInput(
			input,
			"modules",
			"deploy",
			"--module",
			"TST",
			"--semantic-version",
			"1.0.0",
			"--zip-file",
			"test.zip",
			"--platforms",
			"android",
			"--controls",
			"keyboard/mouse",
		)

		Expect(err).To(MatchError("PACKAGE not provided"))
		Expect(output).To(ContainSubstring("PACKAGE"))
	})

	It("can return an error if zip file path is missing", func() {
		input := bytes.NewBufferString("")

		output, err := executor.RunCommandWithInput(
			input,
			"modules",
			"deploy",
			"--module",
			"TST",
			"--semantic-version",
			"1.0.0",
			"--package",
			"test",
			"--platforms",
			"android",
			"--controls",
			"keyboard/mouse",
			"--zip-file",
			"",
		)

		Expect(err).To(MatchError("ZIP FILE not provided"))
		Expect(output).To(ContainSubstring("Enter ZIP FILE:"))
	})

	It("can return an error if the platform options cant be found", func() {
		executor.MockPlatformClient.GetPlatformsError = errors.New("get platforms error")
		input := bytes.NewBufferString("")

		_, err := executor.RunCommandWithInput(
			input,
			"modules",
			"deploy",
			"--module",
			"TST",
			"--semantic-version",
			"1.0.0",
			"--package",
			"test",
			"--zip-file",
			"test.zip",
		)

		Expect(err).To(MatchError("get platforms error"))
		Expect(executor.MockPlatformClient.NumCalledGetPlatforms).To(Equal(1))
	})

	It("can return an error if platforms are missing", func() {
		input := bytes.NewBufferString("")

		output, err := executor.RunCommandWithInput(
			input,
			"modules",
			"deploy",
			"--module",
			"TST",
			"--semantic-version",
			"1.0.0",
			"--package",
			"test",
			"--zip-file",
			"test.zip",
			"--controls",
			"keyboard/mouse",
		)

		Expect(err).To(MatchError("PLATFORMS not provided"))
		Expect(executor.MockPlatformClient.NumCalledGetPlatforms).To(Equal(1))
		Expect(output).To(ContainSubstring("PLATFORMS"))
	})

	It("can return an error if the control types options cant be found", func() {
		executor.MockPlatformClient.GetControlTypesError = errors.New("get controls error")
		input := bytes.NewBufferString("nonexistent\n")

		_, err := executor.RunCommandWithInput(
			input,
			"modules",
			"deploy",
			"--module",
			"TST",
			"--semantic-version",
			"1.0.0",
			"--package",
			"test",
			"--zip-file",
			"test.zip",
			"--platforms",
			"android",
		)

		Expect(err).To(MatchError("get controls error"))
		Expect(executor.MockPlatformClient.NumCalledGetControlTypes).To(Equal(1))
	})

	It("can return an error if control types are missing", func() {
		input := bytes.NewBufferString("\n")

		output, err := executor.RunCommandWithInput(
			input,
			"modules",
			"deploy",
			"--module",
			"TST",
			"--semantic-version",
			"1.0.0",
			"--package",
			"test",
			"--zip-file",
			"test.zip",
			"--platforms",
			"android",
		)

		Expect(err).To(MatchError("CONTROLS not provided"))
		Expect(executor.MockPlatformClient.NumCalledGetControlTypes).To(Equal(1))
		Expect(output).To(ContainSubstring("CONTROLS"))
	})

	It("can return an error if the create call fails", func() {
		executor.MockPlatformClient.CreateModuleVersionError = errors.New("error")
		input := bytes.NewBufferString("\n")

		_, err := executor.RunCommandWithInput(
			input,
			"modules",
			"deploy",
			"--module",
			"TST",
			"--semantic-version",
			"1.0.0",
			"--package",
			"test",
			"--zip-file",
			"test.zip",
			"--platforms",
			"android",
			"--controls",
			"keyboard/mouse",
		)

		Expect(err).To(MatchError("error"))
		Expect(executor.MockPlatformClient.NumCalledCreateModuleVersion).To(Equal(1))
	})

	It("can deploy a module version", func() {
		output, err := executor.RunCommand(
			"modules",
			"deploy",
			"--module",
			"TST",
			"--semantic-version",
			"1.0.0",
			"--package",
			"test",
			"--zip-file",
			"test.zip",
			"--platforms",
			"android",
			"--controls",
			"keyboard/mouse",
		)

		Expect(err).NotTo(HaveOccurred())
		Expect(output).To(ContainSubstring("Deployed version 1.0.0 for module 1"))
		Expect(executor.MockPlatformClient.NumCalledCreateModuleVersion).To(Equal(1))
	})

})
