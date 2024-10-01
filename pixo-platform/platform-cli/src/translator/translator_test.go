package translator_test

import (
	"context"
	"github.com/PixoVR/pixo-golang-clients/pixo-platform/platform-cli/src/translator"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"os"
)

var _ = Describe("Translator", func() {

	var (
		translatorClient *translator.Translator
	)

	BeforeEach(func() {
		var err error
		translatorClient, err = translator.NewTranslator()
		Expect(err).NotTo(HaveOccurred())
		Expect(translatorClient.Model()).To(Equal("icky/translate"))
	})

	It("can translate text", func() {
		if os.Getenv("RUN_INTEGRATION_TESTS") != "true" {
			Skip("Skipping translator test")
		}
		req := translator.Request{
			OriginalLanguage:   "English",
			TranslatedLanguage: "Spanish",
			Text:               "Hello",
		}
		res := make(chan string, 10)
		resFunc := func(resp string) error {
			res <- resp
			return nil
		}
		Expect(translatorClient.TranslateText(context.Background(), req, resFunc)).NotTo(HaveOccurred())
		Expect(<-res).To(Equal("Hola"))
	})

})
