package abstract_client_test

import (
	"fmt"
	"github.com/PixoVR/pixo-golang-clients/pixo-platform/abstract-client"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Abstract", func() {

	var (
		fakeToken = "fake-token"
		apiClient *abstract_client.PixoAbstractAPIClient
	)

	BeforeEach(func() {
		apiClient = abstract_client.NewClient(fakeToken, "")
	})

	It("can set the token", func() {
		newClient := abstract_client.NewClient("", "")
		newClient.SetToken(fakeToken)
		Expect(newClient.GetToken()).To(Equal(fakeToken))
	})

	It("can add headers needed for authentication", func() {
		apiClient.AddHeader("x-fake-header", fakeToken)

		request := apiClient.FormatRequest()

		Expect(request.Header.Get("Authorization")).To(Equal(fmt.Sprintf("Bearer %s", fakeToken)))
		Expect(request.Header.Get("x-fake-header")).To(Equal(fakeToken))
	})

	It("can use the api key", func() {
		apiClient.SetAPIKey(fakeToken)

		request := apiClient.FormatRequest()

		Expect(request.Header.Get("x-api-key")).To(Equal(fakeToken))
	})

	It("should return a response if the request fails", func() {
		res, err := apiClient.Post("invalid", nil)

		Expect(err).To(HaveOccurred())
		Expect(res).To(BeNil())
	})

})
