package parser_test

import (
	"github.com/PixoVR/pixo-golang-clients/pixo-platform/cmd/platform-cli/parser"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Parser", func() {

	var (
		iniParser    *parser.IniParser
		testFilepath = "test-files/test.ini"
	)

	BeforeEach(func() {
		var err error
		iniParser, err = parser.NewIniParser(&testFilepath)
		Expect(err).NotTo(HaveOccurred())
		Expect(iniParser).NotTo(BeNil())
	})

	It("can use the default config file if no config file path is provided", func() {
		iniParser, err := parser.NewIniParser(nil)
		Expect(err).NotTo(HaveOccurred())
		Expect(iniParser).NotTo(BeNil())

		expectedVersion := "3.04.05"
		version, err := iniParser.ParseServerVersion()

		Expect(err).NotTo(HaveOccurred())
		Expect(version).To(Equal(expectedVersion))
	})

	It("can parse the server version from a specific .ini file", func() {
		expectedVersion := "1.02.03"
		version, err := iniParser.ParseServerVersion()

		Expect(err).NotTo(HaveOccurred())
		Expect(version).To(Equal(expectedVersion))
	})

	It("can return an error if the filepath is does not contain a .ini", func() {
		nonexistentFilepath := "test-files/nonexistent.txt"
		_, err := parser.NewIniParser(&nonexistentFilepath)
		Expect(err).To(HaveOccurred())
	})

})
