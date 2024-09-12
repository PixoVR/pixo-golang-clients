package abstract_client_test

import (
	"fmt"
	. "github.com/PixoVR/pixo-golang-clients/pixo-platform/abstract-client"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"net/http"
)

var _ = Describe("Abstract", func() {

	var (
		fakeToken = "fake-token"
		apiClient *AbstractServiceClient
	)

	BeforeEach(func() {
		config := AbstractConfig{
			Token: fakeToken,
		}
		apiClient = NewClient(config)
	})

	It("can set the token", func() {
		config := AbstractConfig{}
		newClient := NewClient(config)
		newClient.SetToken(fakeToken)
		Expect(newClient.GetToken()).To(Equal(fakeToken))
	})

	It("can add headers needed for authentication", func() {
		apiClient.SetHeader("x-fake-header", fakeToken)

		request := apiClient.NewRequest()

		Expect(request.Header.Get("x-fake-header")).To(Equal(fakeToken))
		Expect(request.Header.Get(AuthorizationHeader)).To(Equal(fmt.Sprintf("Bearer %s", fakeToken)))
	})

	It("can use the api key", func() {
		apiClient.SetAPIKey(fakeToken)
		Expect(apiClient.IsAuthenticated()).To(BeTrue())
		Expect(apiClient.GetAPIKey()).To(Equal(fakeToken))

		request := apiClient.NewRequest()

		Expect(request.Header.Get(APIKeyHeader)).To(Equal(fakeToken))
	})

	It("should return a response if the request fails", func() {
		res, err := apiClient.Post("nonexistent", nil)

		Expect(err).NotTo(HaveOccurred())
		Expect(res).NotTo(BeNil())
		Expect(res.StatusCode()).To(Equal(http.StatusNotFound))
	})

	It("can return the current ip address", func() {
		ip, err := apiClient.GetIPAddress()

		Expect(err).NotTo(HaveOccurred())
		Expect(ip).To(MatchRegexp(`\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}`))
	})

})
