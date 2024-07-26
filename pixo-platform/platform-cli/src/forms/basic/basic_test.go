package basic_test

import (
	"bytes"
	"github.com/PixoVR/pixo-golang-clients/pixo-platform/platform-cli/src/forms"
	"github.com/PixoVR/pixo-golang-clients/pixo-platform/platform-cli/src/forms/basic"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"strings"
)

var _ = Describe("Basic Forms", func() {

	var (
		s       *basic.Handler
		input   *bytes.Buffer
		output  *bytes.Buffer
		prompt  = "Do you like pie?"
		options = []forms.Option{
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

	Context("when asking for input", func() {

		It("can ask for input from the user", func() {
			input.WriteString("yes\n")
			question := &forms.Question{
				Type:   forms.Input,
				Prompt: prompt,
			}

			Expect(s.GetResponseFromUser(question)).To(Succeed())

			Expect(question.Answer).NotTo(BeEmpty())
			Expect(question.Answer).To(Equal("yes"))
		})

		It("can ask for sensitive input", func() {
			input.WriteString("password\n")
			question := &forms.Question{Prompt: "Enter password:"}

			Expect(s.GetSensitiveResponseFromUser(question)).To(Succeed())

			Expect(question.Answer).To(Equal("password"))
			Expect(output.String()).To(ContainSubstring("Enter password:"))
			Expect(output.String()).NotTo(ContainSubstring("new-password"))
		})

		It("can read a sensitive value from the user after reading two non sensitive value", func() {
			input.WriteString("new-username\nnewer-username\nnew-password\n")
			question := &forms.Question{
				Type:   forms.Input,
				Prompt: "username",
			}

			Expect(s.GetResponseFromUser(question)).To(Succeed())
			Expect(question.Answer).To(Equal("new-username"))

			Expect(s.GetResponseFromUser(question)).To(Succeed())
			Expect(question.Answer).To(Equal("newer-username"))

			question.Type = forms.SensitiveInput
			Expect(s.GetSensitiveResponseFromUser(question)).To(Succeed())
			Expect(question.Answer).To(Equal("new-password"))
		})

	})

	DescribeTable("can ask for confirmation",
		func(input string, expected bool) {
			question := &forms.Question{
				Type:   forms.Confirm,
				Prompt: prompt,
			}
			inputBuffer := bytes.NewBufferString(input)
			s.SetReader(inputBuffer)

			Expect(s.Confirm(question)).To(Succeed())
			Expect(question.Answer).To(Equal(expected))
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
			question *forms.Question
		)

		BeforeEach(func() {
			question = &forms.Question{
				Type:    forms.Select,
				Prompt:  "Select one",
				Options: options,
			}
		})

		It("can return an error if no option is selected", func() {
			input.WriteString("\n")

			err := s.Select(question)

			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("not provided"))
			Expect(output.String()).To(ContainSubstring(question.Prompt))
			for _, option := range question.Options {
				Expect(output.String()).To(ContainSubstring(option.Label))
			}
		})

		It("can return an error if the option is not valid", func() {
			input.WriteString("probably\n")

			err := s.Select(question)

			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("invalid option"))
		})

		It("can ask a single select question", func() {
			input.WriteString("yes\n")
			Expect(s.Select(question)).To(Succeed())
			Expect(question.Answer).To(Equal("yes"))
		})

		It("can retrieve the options from a function", func() {
			input.WriteString("functional-no\n")
			question.Options = nil
			question.GetOptionsFunc = func() ([]forms.Option, error) {
				options := []forms.Option{
					{Label: "functional-yes", Value: "1"},
					{Label: "functional-no", Value: "2"},
				}
				return options, nil
			}

			Expect(s.Select(question)).To(Succeed())

			Expect(question.Answer).To(Equal("functional-no"))
		})

		It("can return an error if no option is selected when getting ids", func() {
			input.WriteString("\n")

			err := s.SelectID(question)

			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("not provided"))
			Expect(output.String()).To(ContainSubstring(question.Prompt))
			for _, option := range question.Options {
				Expect(output.String()).To(ContainSubstring(option.Label))
			}
		})

		It("can ask a single select question with ids", func() {
			input.WriteString("yes\n")
			Expect(s.SelectID(question)).To(Succeed())
			Expect(question.Answer).To(Equal(1))
		})

		It("can retrieve the options for a function with ids", func() {
			input.WriteString("functional-yes\n")
			question.Options = nil
			question.GetOptionsFunc = func() ([]forms.Option, error) {
				options := []forms.Option{
					{Label: "functional-yes", Value: "1"},
					{Label: "functional-no", Value: "2"},
				}
				return options, nil
			}

			Expect(s.SelectID(question)).To(Succeed())

			Expect(question.Answer).To(Equal(1))
		})

	})

	Context("multiselect", func() {

		var (
			question *forms.Question
		)

		BeforeEach(func() {
			question = &forms.Question{
				Type:    forms.MultiSelect,
				Prompt:  "Select multiple",
				Options: options,
			}
		})

		It("can return an error if one of the options is invalid", func() {
			input.WriteString("yes,probably\n")

			err := s.MultiSelect(question)

			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("invalid option"))
		})

		It("can ask a question with multiple answers", func() {
			input.WriteString("yes\n")

			Expect(s.MultiSelect(question)).To(Succeed())

			Expect(question.Answer).To(HaveLen(1))
			Expect(question.Answer.([]string)[0]).To(Equal("yes"))
			Expect(output.String()).To(ContainSubstring(question.Prompt))
			for _, option := range question.Options {
				Expect(output.String()).To(ContainSubstring(option.Label))
			}
		})

		It("can retrieve the options from a function", func() {
			input.WriteString("functional-no\n")
			question.Options = nil
			question.GetOptionsFunc = func() ([]forms.Option, error) {
				options := []forms.Option{
					{Label: "functional-yes", Value: "1"},
					{Label: "functional-no", Value: "2"},
				}
				return options, nil
			}

			Expect(s.MultiSelect(question)).To(Succeed())

			Expect(question.Answer).To(HaveLen(1))
			Expect(question.Answer.([]string)[0]).To(Equal("functional-no"))
		})

		It("can ask a question with multiple answers and return values as ints", func() {
			input.WriteString("yes,no\n")

			Expect(s.MultiSelectIDs(question)).To(Succeed())

			Expect(question.Answer).To(HaveLen(2))
			Expect(question.Answer.([]int)[0]).To(Equal(1))
			Expect(question.Answer.([]int)[1]).To(Equal(2))
			Expect(output.String()).To(ContainSubstring(question.Prompt))
			for _, option := range question.Options {
				Expect(output.String()).To(ContainSubstring(option.Label))
			}
		})

		It("can ask a question with multiple answers with overriding io", func() {
			customInput := &bytes.Buffer{}
			customOutput := &bytes.Buffer{}
			s.SetReader(customInput)
			s.SetWriter(customOutput)
			customInput.WriteString("yes,no\n")

			Expect(s.MultiSelect(question)).To(Succeed())

			Expect(question.Answer).To(HaveLen(2))
			Expect(question.Answer.([]string)[0]).To(Equal("yes"))
			Expect(question.Answer.([]string)[1]).To(Equal("no"))
			Expect(customOutput.String()).To(ContainSubstring(question.Prompt))
			Expect(customOutput.String()).To(ContainSubstring("yes"))
			Expect(customOutput.String()).To(ContainSubstring("no"))
		})

		It("can retrieve the options from a function with ids", func() {
			input.WriteString("functional-no\n")
			question.Options = nil
			question.GetOptionsFunc = func() ([]forms.Option, error) {
				options := []forms.Option{
					{Label: "functional-yes", Value: "1"},
					{Label: "functional-no", Value: "2"},
				}
				return options, nil
			}

			Expect(s.MultiSelectIDs(question)).To(Succeed())

			Expect(question.Answer).To(HaveLen(1))
			Expect(question.Answer.([]int)[0]).To(Equal(2))
		})

		It("can display an entire form", func() {
			inputLines := []string{
				"some-input",
				"sensitive-input",
				"no",
				"yes,no",
				"no,yes",
				"one",
				"two",
			}
			input.WriteString(strings.Join(inputLines, "\n") + "\n")

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
				{
					Key:      "optional-input",
					Prompt:   "Enter some optional input",
					Type:     forms.Input,
					Optional: true,
				},
			}

			answers, err := s.AskQuestions(questions)

			Expect(err).NotTo(HaveOccurred())
			Expect(answers).NotTo(BeNil())
			Expect(answers).To(HaveLen(len(questions) - 1))

			Expect(forms.String(answers["something"])).To(Equal(inputLines[0]))
			Expect(output.String()).To(ContainSubstring("Enter some input"))

			Expect(forms.String(answers["sensitive"])).To(Equal(inputLines[1]))
			Expect(output.String()).To(ContainSubstring("Enter sensitive output"))
			Expect(output.String()).NotTo(ContainSubstring("sensitive-input"))

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

			Expect(answers["optional-select"]).To(BeNil())
			Expect(output.String()).To(ContainSubstring("Enter some optional input"))
		})

		//It("can run a function to retrieve the data if not provided", func() {
		//})

	})

})
