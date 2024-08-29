package platform_test

import (
	"context"
	. "github.com/PixoVR/pixo-golang-clients/pixo-platform/platform"
	"github.com/go-faker/faker/v4"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Webhooks", func() {

	var (
		ctx          context.Context
		webhookInput Webhook
		testWebhook  *Webhook
	)

	BeforeEach(func() {
		ctx = context.Background()

		webhookInput = Webhook{
			OrgID:         1,
			URL:           "http://example.com",
			EventTypes:    []string{"SessionCompleted"},
			GenerateToken: &[]bool{true}[0],
			Description:   faker.Sentence(),
		}
		var err error
		testWebhook, err = tokenClient.CreateWebhook(ctx, webhookInput)
		Expect(err).NotTo(HaveOccurred())
		Expect(testWebhook).NotTo(BeNil())
		Expect(testWebhook.ID).NotTo(BeZero())
		Expect(testWebhook.OrgID).To(Equal(webhookInput.OrgID))
		Expect(testWebhook.Org).NotTo(BeNil())
		Expect(testWebhook.Org.ID).To(Equal(webhookInput.OrgID))
		Expect(testWebhook.URL).To(Equal(webhookInput.URL))
		Expect(testWebhook.Description).To(Equal(webhookInput.Description))
		Expect(testWebhook.Token).NotTo(BeEmpty())
	})

	AfterEach(func() {
		err := tokenClient.DeleteWebhook(ctx, testWebhook.ID)
		Expect(err).NotTo(HaveOccurred())
		deletedWebhook, err := tokenClient.GetWebhook(ctx, webhookInput.ID)
		Expect(err).To(HaveOccurred())
		Expect(deletedWebhook).To(BeNil())
	})

	It("can get webhooks by org id", func() {
		webhooks, err := tokenClient.GetWebhooks(ctx, &WebhookParams{OrgID: testWebhook.OrgID})
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
		webhookInput.ID = testWebhook.ID
		webhookInput.URL = "http://example.com/updated"
		webhookInput.Description = "updated description"
		webhookInput.GenerateToken = &[]bool{false}[0]
		webhookInput.Token = "updated-token"
		webhookInput.EventTypes = []string{"UserCreated"}

		updatedWebhook, err := tokenClient.UpdateWebhook(ctx, webhookInput)

		Expect(err).NotTo(HaveOccurred())
		Expect(updatedWebhook).NotTo(BeNil())
		Expect(updatedWebhook.ID).To(Equal(webhookInput.ID))
		Expect(updatedWebhook.URL).To(Equal(webhookInput.URL))
		Expect(updatedWebhook.Description).To(Equal(webhookInput.Description))
		Expect(updatedWebhook.Token).To(Equal(webhookInput.Token))
	})

})
