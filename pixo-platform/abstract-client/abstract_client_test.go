package abstract_client_test

import (
	"github.com/PixoVR/pixo-golang-clients/pixo-platform/abstract-client"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"net/http"
)

var _ = Describe("Abstract", func() {

	It("can set the token", func() {
		client := abstract_client.NewClient("", "")
		client.SetToken("some-fake-token")
		Expect(client.GetToken()).To(Equal("some-fake-token"))
	})

	It("can format the request headers needed for authentication", func() {
		client := abstract_client.NewClient("some-fake-token", "")
		client.AddHeader("x-fake-header", "some-fake-token")

		request := client.FormatRequest()

		Expect(request.Header.Get("Authorization")).To(Equal("Bearer some-fake-token"))
		Expect(request.Header.Get("x-access-token")).To(Equal("some-fake-token"))
		Expect(request.Header.Get("x-fake-header")).To(Equal("some-fake-token"))
	})

	It("should be able to add headers to the request", func() {
		client := abstract_client.NewClient("", "https://api.apex.dev.pixovr.com")

		res, err := client.Get("health")

		Expect(err).NotTo(HaveOccurred())
		Expect(res.StatusCode()).To(Equal(http.StatusOK))
	})

	It("should return a response if the request fails", func() {
		client := abstract_client.NewClient("", "")

		res, err := client.Post("invalid", nil)

		Expect(err).To(HaveOccurred())
		Expect(res).To(BeNil())
	})

})
