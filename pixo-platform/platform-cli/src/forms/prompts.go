package forms

import "strings"

func CleanPrompt(prompt string) (cleanPrompt string) {
	cleanPrompt = strings.ReplaceAll(prompt, "-", " ")
	cleanPrompt = strings.ToUpper(cleanPrompt)

	cleanPrompt = strings.ReplaceAll(cleanPrompt, " IDS ", " IDs ")
	if strings.HasSuffix(cleanPrompt, " IDS") {
		cleanPrompt = strings.TrimSuffix(cleanPrompt, " IDS")
		cleanPrompt = cleanPrompt + " IDs"
	}
	return cleanPrompt
}
