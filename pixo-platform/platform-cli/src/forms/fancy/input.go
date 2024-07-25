package fancy

import "github.com/charmbracelet/huh"

func (f *Handler) InputField(prompt string, response *string) huh.Field {
	return huh.NewInput().
		Title(prompt).
		Value(response)
}

func (f *Handler) GetResponseFromUser(prompt string) (string, error) {
	var response string
	if err := f.InputField(prompt, &response).Run(); err != nil {
		return "", err
	}

	return response, nil
}

func (f *Handler) SensitiveInputField(prompt string, response *string) huh.Field {
	return huh.NewInput().
		Title(prompt).
		EchoMode(huh.EchoModePassword).
		Value(response)
}

func (f *Handler) GetSensitiveResponseFromUser(prompt string) (string, error) {
	var response string
	if err := f.SensitiveInputField(prompt, &response).Run(); err != nil {
		return "", err
	}

	return response, nil
}
