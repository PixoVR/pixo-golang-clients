package fancy

import "github.com/charmbracelet/huh"

func (f *Handler) InputField(prompt string, response *string) huh.Field {
	return huh.NewInput().
		Title(prompt).
		Value(response)
}

func (f *Handler) GetResponseFromUser(prompt string, response *string) error {
	return f.InputField(prompt, response).Run()
}

func (f *Handler) SensitiveInputField(prompt string, response *string) huh.Field {
	return huh.NewInput().
		Title(prompt).
		EchoMode(huh.EchoModePassword).
		Value(response)
}

func (f *Handler) GetSensitiveResponseFromUser(prompt string, response *string) error {
	return f.SensitiveInputField(prompt, response).Run()
}
