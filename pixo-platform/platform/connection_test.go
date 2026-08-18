package platform_test

import (
	"context"

	. "github.com/PixoVR/pixo-golang-clients/pixo-platform/platform"
	"github.com/PixoVR/pixo-golang-clients/pixo-platform/urlfinder"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Connection Check", func() {

	It("can report when no credentials are configured", func() {
		client := NewClient(urlfinder.ClientConfig{Lifecycle: lifecycle})
		client.SetAPIKey("")

		err := client.CheckConnection(context.Background())

		Expect(err).To(MatchError(ErrNoCredentials))
		Expect(err.Error()).To(ContainSubstring(client.GetURL()))
	})

	It("can report when the api key is rejected", func() {
		client := NewClient(urlfinder.ClientConfig{Lifecycle: lifecycle, APIKey: "not-a-real-api-key"})

		err := client.CheckConnection(context.Background())

		Expect(err).To(MatchError(ErrUnauthorized))
		Expect(err.Error()).To(ContainSubstring("api key"))
		Expect(err.Error()).To(ContainSubstring(client.GetURL()))
	})

	It("can report when the token is rejected", func() {
		client := NewClient(urlfinder.ClientConfig{Lifecycle: lifecycle, Token: "not-a-real-token"})

		err := client.CheckConnection(context.Background())

		Expect(err).To(MatchError(ErrUnauthorized))
		Expect(err.Error()).To(ContainSubstring("token"))
	})

	It("can verify a working api key", func() {
		Expect(apiKeyClient.CheckConnection(context.Background())).To(Succeed())
	})

	It("can verify a working token", func() {
		Expect(tokenClient.CheckConnection(context.Background())).To(Succeed())
	})

})
