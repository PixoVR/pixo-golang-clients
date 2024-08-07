package legacy_test

import (
	primary_api "github.com/PixoVR/pixo-golang-clients/pixo-platform/legacy"
	"github.com/PixoVR/pixo-golang-clients/pixo-platform/urlfinder"
	config2 "github.com/PixoVR/pixo-golang-server-utilities/pixo-platform/config"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"os"
)

var _ = Describe("Webhook", Ordered, func() {

	var (
		primaryAPIClient *primary_api.Client
		webhook          = primary_api.Webhook{
			OrgID:       20,
			Description: "test-webhook",
			URL:         "https://example.com",
		}
	)

	BeforeEach(func() {
		config := urlfinder.ClientConfig{
			Lifecycle: config2.GetEnvOrReturn("TEST_PIXO_LIFECYCLE", "dev"),
			Region:    config2.GetEnvOrReturn("TEST_PIXO_REGION", "na"),
		}
		primaryAPIClient = primary_api.NewClient(config)
		Expect(primaryAPIClient.Login(os.Getenv("TEST_PIXO_USERNAME"), os.Getenv("TEST_PIXO_PASSWORD"))).To(Succeed())
	})

	It("can create a webhook", func() {
		err := primaryAPIClient.CreateWebhook(webhook)

		Expect(err).NotTo(HaveOccurred())
	})

	It("can get webhooks", func() {
		webhooks, err := primaryAPIClient.GetWebhooks(webhook.OrgID)

		Expect(err).NotTo(HaveOccurred())
		Expect(webhooks).NotTo(BeNil())
		Expect(len(webhooks)).To(BeNumerically(">", 0))
	})

	It("can delete a webhook", func() {
		webhooks, err := primaryAPIClient.GetWebhooks(webhook.OrgID)
		Expect(err).NotTo(HaveOccurred())
		Expect(webhooks).NotTo(BeNil())
		Expect(len(webhooks)).To(BeNumerically(">", 0))

		for _, w := range webhooks {
			if w.Description == "test-webhook" {
				err = primaryAPIClient.DeleteWebhook(w.ID)
				Expect(err).NotTo(HaveOccurred())
			}
		}
	})

})
