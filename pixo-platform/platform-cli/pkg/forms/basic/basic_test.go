package basic_test

import (
	"bytes"
	"github.com/PixoVR/pixo-golang-clients/pixo-platform/platform-cli/pkg/forms"
	"github.com/PixoVR/pixo-golang-clients/pixo-platform/platform-cli/pkg/forms/basic"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Basic Forms", func() {

	var (
		s        *basic.FormHandler
		input    *bytes.Buffer
		output   *bytes.Buffer
		question = "Do you like pie?"
		options  = []forms.Option{
			{Label: "yes", Value: "1"},
			{Label: "no", Value: "2"},
			{Label: "maybe", Value: "3"},
		}
	)

	BeforeEach(func() {
		input = &bytes.Buffer{}
		output = &bytes.Buffer{}
		s = basic.NewFormHandler(input, output)
		Expect(s).NotTo(BeNil())
	})

	It("can ask for basic input", func() {
		input.WriteString("yes\n")

		answer, err := s.GetResponseFromUser(question)

		Expect(err).NotTo(HaveOccurred())
		Expect(answer).NotTo(BeEmpty())
		Expect(answer).To(Equal("yes"))
	})

	It("can ask a question with multiple answers", func() {
		input.WriteString("yes\n")

		answers, err := s.MultiSelect(question, options)

		Expect(err).NotTo(HaveOccurred())
		Expect(answers).To(HaveLen(1))
		Expect(answers[0]).To(Equal("yes"))
		Expect(output.String()).To(ContainSubstring(question))
		for _, option := range options {
			Expect(output.String()).To(ContainSubstring(option.Label))
		}
	})

	It("can ask a question with multiple answers and return values as ints", func() {
		input.WriteString("yes,no\n")

		answers, err := s.MultiSelectIDs(question, options)

		Expect(err).NotTo(HaveOccurred())
		Expect(answers).To(HaveLen(2))
		Expect(answers[0]).To(Equal(1))
		Expect(answers[1]).To(Equal(2))
		Expect(output.String()).To(ContainSubstring(question))
	})

	It("can ask a question with multiple answers with custom io", func() {
		customInput := &bytes.Buffer{}
		customOutput := &bytes.Buffer{}
		s.SetReader(customInput)
		s.SetWriter(customOutput)

		customInput.WriteString("yes,no\n")

		answers, err := s.MultiSelect(question, options)

		Expect(err).NotTo(HaveOccurred())
		Expect(answers).To(HaveLen(2))
		Expect(answers[0]).To(Equal("yes"))
		Expect(answers[1]).To(Equal("no"))
		Expect(customOutput.String()).To(ContainSubstring(question))
		Expect(customOutput.String()).To(ContainSubstring("yes"))
		Expect(customOutput.String()).To(ContainSubstring("no"))
	})

})
