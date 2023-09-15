package abstract_client_test

import (
	"github.com/PixoVR/pixo-golang-clients/pixo-platform/abstract-client"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"net/http"
	"os"
)

var _ = Describe("Abstract", func() {

	It("should return the correct URL", func() {
		client := abstract_client.NewClient(os.Getenv("API_TOKEN"), "")
		Expect(client.GetURL()).To(Equal("https://api.apex.pixovr.com"))
	})

	It("can set the token", func() {
		client := abstract_client.NewClient("", "")
		client.SetToken("some-fake-token")
		Expect(client.GetToken()).To(Equal("some-fake-token"))
	})

	It("can format the request headers needed for authentication", func() {
		client := abstract_client.NewClient("some-fake-token", "")

		request := client.FormatRequest()
		Expect(request.Header.Get("x-access-token")).To(Equal("some-fake-token"))
		Expect(request.Header.Get("Authorization")).To(Equal("Bearer some-fake-token"))
	})

	It("should be able to make a get request", func() {
		client := abstract_client.NewClient("", "")
		res, err := client.Get("health")
		Expect(err).NotTo(HaveOccurred())
		Expect(res.StatusCode()).To(Equal(http.StatusOK))
	})

})
