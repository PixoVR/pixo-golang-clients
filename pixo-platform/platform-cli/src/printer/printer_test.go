package printer_test

import (
	"bytes"
	"fmt"
	. "github.com/PixoVR/pixo-golang-clients/pixo-platform/platform-cli/src/printer"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Emoji Printer", func() {

	var (
		printer *EmojiPrinter
		input   *bytes.Buffer
		output  *bytes.Buffer
	)

	BeforeEach(func() {
		input = &bytes.Buffer{}
		output = &bytes.Buffer{}
		printer = NewEmojiPrinter(output)
	})

	AfterEach(func() {
		input.Reset()
		output.Reset()
	})

	It("can output a messages with emojis", func() {
		msg := "hello world\n"

		printer.Print(msg)
		Expect(output.String()).To(Equal(msg))
		output.Reset()

		printer.Println(msg)
		Expect(output.String()).To(Equal(msg + "\n"))
		output.Reset()

		msg = fmt.Sprintf(":rocket:hello %s", "world")
		expectedMsg := "🚀 hello world"
		printer.Print(msg)
		Expect(output.String()).To(Equal(expectedMsg))
	})

	It("can override the writer", func() {
		outputWriter := bytes.NewBufferString("")
		printer.SetWriter(outputWriter)
		msg := "hello world\n"

		printer.Print(msg)
		Expect(outputWriter.String()).To(Equal(msg))
	})

})
