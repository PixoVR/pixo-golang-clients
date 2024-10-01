package translator

type Request struct {
	OriginalLanguage   string
	TranslatedLanguage string
	Type               string
	Text               string
	Filepath           string
}
