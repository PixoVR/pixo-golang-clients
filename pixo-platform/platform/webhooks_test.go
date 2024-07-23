package platform_test

import (
	"context"
	"github.com/PixoVR/pixo-golang-clients/pixo-platform/platform"
	. "github.com/PixoVR/pixo-golang-clients/pixo-platform/platform"
	"github.com/go-faker/faker/v4"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Webhooks", func() {

	var (
		ctx          context.Context
		webhookInput platform.Webhook
		testWebhook  *platform.Webhook
	)

	BeforeEach(func() {
		ctx = context.Background()

		webhookInput = platform.Webhook{
			OrgID:       1,
			URL:         "http://example.com",
			Token:       "token",
			Description: faker.Sentence(),
		}
		var err error
		testWebhook, err = tokenClient.CreateWebhook(ctx, webhookInput)
		Expect(err).NotTo(HaveOccurred())
		Expect(testWebhook).NotTo(BeNil())
		Expect(testWebhook.ID).NotTo(BeZero())
	})

	AfterEach(func() {
		err := tokenClient.DeleteWebhook(ctx, testWebhook.ID)
		Expect(err).NotTo(HaveOccurred())
		deletedWebhook, err := tokenClient.GetWebhook(ctx, webhookInput.ID)
		Expect(err).To(HaveOccurred())
		Expect(deletedWebhook).To(BeNil())
	})

	It("can get webhooks by org id", func() {
		webhooks, err := tokenClient.GetWebhooks(ctx, WebhookParams{OrgID: testWebhook.OrgID})
		Expect(err).NotTo(HaveOccurred())
		Expect(webhooks).NotTo(BeNil())
		Expect(len(webhooks)).To(BeNumerically(">", 0))
	})

	It("can return an error if webhook does not exist when updating", func() {
		updatedWebhook, err := tokenClient.UpdateWebhook(ctx, webhookInput)
		Expect(err).To(HaveOccurred())
		Expect(updatedWebhook).To(BeNil())
		Expect(err.Error()).To(ContainSubstring("webhook id is required"))
	})

	It("can update a webhook", func() {
		webhookToUpdateInput := Webhook{
			OrgID:       1,
			URL:         "http://example.com",
			Description: "client test",
			Token:       "token",
		}
		webhookToUpdate, err := tokenClient.CreateWebhook(ctx, webhookToUpdateInput)
		Expect(err).NotTo(HaveOccurred())
		Expect(webhookToUpdate).NotTo(BeNil())
		Expect(webhookToUpdate.ID).NotTo(BeZero())
		webhookToUpdateInput.ID = webhookToUpdate.ID
		webhookToUpdateInput.URL = "http://example.com/updated"
		webhookToUpdateInput.Description = "updated description"

		updatedWebhook, err := tokenClient.UpdateWebhook(ctx, webhookToUpdateInput)

		Expect(err).NotTo(HaveOccurred())
		Expect(updatedWebhook).NotTo(BeNil())
		Expect(updatedWebhook.ID).To(Equal(webhookToUpdateInput.ID))
		Expect(updatedWebhook.URL).To(Equal(webhookToUpdateInput.URL))
		Expect(updatedWebhook.Description).To(Equal(webhookToUpdateInput.Description))
		Expect(updatedWebhook.Token).To(Equal(webhookToUpdateInput.Token))
	})

})
