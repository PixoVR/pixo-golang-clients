package vertexai_test

import (
	vertexai "github.com/PixoVR/pixo-golang-clients/vertex-ai/client"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Vertexai", func() {

	It("can get a comment using vertex AI code model", func() {
		client, err := vertexai.NewClient(vertexai.ClientConfig{})
		Expect(err).NotTo(HaveOccurred())

		req := vertexai.ChatRequest{
			ModelID: vertexai.BISON_CODE_MODEL_ID,
			Prompt:  "Who is the best basketball player of all time?",
		}
		response, err := client.ChatResponse(req)
		Expect(err).NotTo(HaveOccurred())
		Expect(response).NotTo(BeEmpty())
	})

	It("can get a comment using vertex AI text model", func() {
		client, err := vertexai.NewClient(vertexai.ClientConfig{Debug: true})
		Expect(err).NotTo(HaveOccurred())

		req := vertexai.ChatRequest{
			ModelID: vertexai.BISON_TEXT_MODEL_ID,
			Instances: []vertexai.Instance{
				{
					Content: "Q: What do you call a fake noodle?",
				},
			},
		}
		response, err := client.ChatResponse(req)
		Expect(err).NotTo(HaveOccurred())
		Expect(response).NotTo(BeEmpty())
	})

})
