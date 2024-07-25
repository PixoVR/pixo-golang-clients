package fancy

import "github.com/charmbracelet/huh"

func (f *Handler) ConfirmField(title string, confirm *bool) huh.Field {
	return huh.NewConfirm().
		Title(title).
		Affirmative("Yes").
		Negative("No").
		Value(confirm)
}

func (f *Handler) Confirm(title string, confirm *bool) error {
	return f.ConfirmField(title, confirm).Run()
}
