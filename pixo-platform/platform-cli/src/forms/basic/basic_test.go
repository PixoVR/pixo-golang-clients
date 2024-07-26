package basic_test

import (
	"bytes"
	"github.com/PixoVR/pixo-golang-clients/pixo-platform/platform-cli/src/forms"
	"github.com/PixoVR/pixo-golang-clients/pixo-platform/platform-cli/src/forms/basic"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Basic Forms", func() {

	var (
		s        *basic.Handler
		input    *bytes.Buffer
		output   *bytes.Buffer
		question = "Do you like pie?"
		options  = []forms.Option{
			{Label: "yes", Value: "1"},
			{Label: "no", Value: "2"},
			{Label: "maybe", Value: "3"},
		}
		answer string
	)

	BeforeEach(func() {
		input = &bytes.Buffer{}
		output = &bytes.Buffer{}
		s = basic.NewFormHandler(input, output)
		Expect(s).NotTo(BeNil())
	})

	Context("when asking for input", func() {

		It("can ask for basic input", func() {
			input.WriteString("yes\n")

			err := s.GetResponseFromUser(question, &answer)

			Expect(err).NotTo(HaveOccurred())
			Expect(answer).NotTo(BeEmpty())
			Expect(answer).To(Equal("yes"))
		})

		It("can ask for sensitive input", func() {
			input.WriteString("password\n")

			err := s.GetSensitiveResponseFromUser("Enter password:", &answer)

			Expect(err).NotTo(HaveOccurred())
			Expect(answer).To(Equal("password"))
			Expect(output.String()).To(ContainSubstring("Enter password:"))
			Expect(output.String()).NotTo(ContainSubstring("new-password"))
		})

		It("can read a sensitive value from the user after reading a non sensitive value", func() {
			input.WriteString("new-username\nnewer-username\nnew-password\n")

			err := s.GetResponseFromUser("username", &answer)
			Expect(err).NotTo(HaveOccurred())
			Expect(answer).To(Equal("new-username"))

			err = s.GetResponseFromUser("username", &answer)
			Expect(err).NotTo(HaveOccurred())
			Expect(answer).To(Equal("newer-username"))

			err = s.GetSensitiveResponseFromUser("password", &answer)
			Expect(err).NotTo(HaveOccurred())
			Expect(answer).To(Equal("new-password"))
		})

	})

	DescribeTable("can ask for confirmation",
		func(input string, expected bool) {
			var boolAnswer bool
			inputBuffer := bytes.NewBufferString(input)
			s.SetReader(inputBuffer)

			err := s.Confirm(question, &boolAnswer)
			Expect(err).NotTo(HaveOccurred())

			Expect(boolAnswer).To(Equal(expected))
		},
		Entry("empty", "\n", false),
		Entry("yes", "yes\n", true),
		Entry("YES", "YES\n", true),
		Entry("y", "y\n", true),
		Entry("Y", "Y\n", true),
		Entry("no", "no\n", false),
		Entry("NO", "NO\n", false),
		Entry("n", "n\n", false),
		Entry("N", "N\n", false),
	)

	Context("single select", func() {

		var (
			intID int
		)

		It("can return an error if no option is selected", func() {
			input.WriteString("\n")

			err := s.Select(question, options, &answer)

			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("not provided"))
			Expect(output.String()).To(ContainSubstring(question))
			for _, option := range options {
				Expect(output.String()).To(ContainSubstring(option.Label))
			}
		})

		It("can ask a single select question", func() {
			input.WriteString("yes\n")
			Expect(s.Select(question, options, &answer)).To(Succeed())
			Expect(answer).To(Equal("yes"))
		})

		It("can return an error if no option is selected when getting ids", func() {
			input.WriteString("\n")

			err := s.SelectID(question, options, &intID)

			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("not provided"))
			Expect(output.String()).To(ContainSubstring(question))
			for _, option := range options {
				Expect(output.String()).To(ContainSubstring(option.Label))
			}
		})

		It("can ask a single select question with ids", func() {
			input.WriteString("yes\n")
			Expect(s.SelectID(question, options, &intID)).To(Succeed())
			Expect(intID).To(Equal(1))
		})

	})

	Context("multiselect", func() {

		var (
			intSliceAnswer    []int
			stringSliceAnswer []string
		)

		It("can ask a question with multiple answers", func() {
			input.WriteString("yes\n")

			err := s.MultiSelect(question, options, &stringSliceAnswer)

			Expect(err).NotTo(HaveOccurred())
			Expect(stringSliceAnswer).To(HaveLen(1))
			Expect(stringSliceAnswer[0]).To(Equal("yes"))
			Expect(output.String()).To(ContainSubstring(question))
			for _, option := range options {
				Expect(output.String()).To(ContainSubstring(option.Label))
			}
		})

		It("can ask a question with multiple answers and return values as ints", func() {
			input.WriteString("yes,no\n")

			err := s.MultiSelectIDs(question, options, &intSliceAnswer)

			Expect(err).NotTo(HaveOccurred())
			Expect(intSliceAnswer).To(HaveLen(2))
			Expect(intSliceAnswer[0]).To(Equal(1))
			Expect(intSliceAnswer[1]).To(Equal(2))
			Expect(output.String()).To(ContainSubstring(question))
		})

		It("can ask a question with multiple answers with overriding io", func() {
			customInput := &bytes.Buffer{}
			customOutput := &bytes.Buffer{}
			s.SetReader(customInput)
			s.SetWriter(customOutput)

			customInput.WriteString("yes,no\n")

			err := s.MultiSelect(question, options, &stringSliceAnswer)

			Expect(err).NotTo(HaveOccurred())
			Expect(stringSliceAnswer).To(HaveLen(2))
			Expect(stringSliceAnswer[0]).To(Equal("yes"))
			Expect(stringSliceAnswer[1]).To(Equal("no"))
			Expect(customOutput.String()).To(ContainSubstring(question))
			Expect(customOutput.String()).To(ContainSubstring("yes"))
			Expect(customOutput.String()).To(ContainSubstring("no"))
		})

		It("can display an entire form", func() {
			input.WriteString("some-input\nsensitive-output\nno\nyes,no\nno,yes\none\ntwo\n")

			questions := []forms.Question{
				{
					Key:    "something",
					Prompt: "Enter some input",
					Type:   forms.Input,
				},
				{
					Key:    "sensitive",
					Prompt: "Enter sensitive output",
					Type:   forms.SensitiveInput,
				},
				{
					Key:    "confirm",
					Prompt: "Do you like pie?",
					Type:   forms.Confirm,
				},
				{
					Key:    "multiselect",
					Prompt: "Select multiple options",
					Type:   forms.MultiSelect,
					Options: []forms.Option{
						{Label: "yes"},
						{Label: "no"},
					},
				},
				{
					Key:    "multiselect-ids",
					Prompt: "Select multiple options with ids",
					Type:   forms.MultiSelectIDs,
					Options: []forms.Option{
						{Label: "yes", Value: "1"},
						{Label: "no", Value: "2"},
					},
				},
				{
					Key:    "select",
					Prompt: "Select",
					Type:   forms.Select,
					Options: []forms.Option{
						{Label: "one"},
						{Label: "two"},
					},
				},
				{
					Key:    "select-id",
					Prompt: "Select ID",
					Type:   forms.SelectID,
					Options: []forms.Option{
						{Label: "one", Value: "1"},
						{Label: "two", Value: "2"},
					},
				},
			}

			answers, err := s.AskQuestions(questions)

			Expect(err).NotTo(HaveOccurred())
			Expect(answers).NotTo(BeNil())

			Expect(forms.String(answers["something"])).To(Equal("some-input"))
			Expect(output.String()).To(ContainSubstring("Enter some input"))

			Expect(forms.String(answers["sensitive"])).To(Equal("sensitive-output"))
			Expect(output.String()).To(ContainSubstring("Enter sensitive output"))
			Expect(output.String()).NotTo(ContainSubstring("sensitive-output"))

			Expect(forms.Bool(answers["confirm"])).To(BeFalse())
			Expect(output.String()).To(ContainSubstring("Enter some input"))

			Expect(forms.StringSlice(answers["multiselect"])).To(Equal([]string{"yes", "no"}))
			Expect(output.String()).To(ContainSubstring("Select multiple options"))

			Expect(forms.IntSlice(answers["multiselect-ids"])).To(Equal([]int{2, 1}))
			Expect(output.String()).To(ContainSubstring("Select multiple options with ids"))

			Expect(forms.String(answers["select"])).To(Equal("one"))
			Expect(output.String()).To(ContainSubstring("Select"))

			Expect(forms.Int(answers["select-id"])).To(Equal(2))
			Expect(output.String()).To(ContainSubstring("Select ID"))
		})

	})

})
