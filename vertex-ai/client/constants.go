package vertexai

const (
	API_ENDPOINT        = "us-central1-aiplatform.googleapis.com"
	BISON_CODE_MODEL_ID = "code-bison@001"
	BISON_TEXT_MODEL_ID = "text-bison@001"
)

var (
	DEFAULT_PARAMETER = Parameters{
		MaxOutputTokens: 256,
		Temperature:     0.2,
		TopP:            .8,
		TopK:            40,
	}
)
