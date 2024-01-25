package graphql_api_test

import (
	"context"
	. "github.com/PixoVR/pixo-golang-clients/pixo-platform/graphql-api"
	platform "github.com/PixoVR/pixo-golang-clients/pixo-platform/primary-api"
	"github.com/PixoVR/pixo-golang-clients/pixo-platform/urlfinder"
	"github.com/go-faker/faker/v4"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Users API", func() {

	var (
		ctx       context.Context
		userInput platform.User
		testUser  *platform.User
	)

	BeforeEach(func() {
		ctx = context.Background()

		userInput = platform.User{
			FirstName: faker.FirstName(),
			LastName:  faker.LastName(),
			Username:  faker.Username(),
			Password:  faker.Password(),
			OrgID:     1,
		}
		var err error
		testUser, err = tokenClient.CreateUser(ctx, userInput)
		Expect(err).NotTo(HaveOccurred())
		Expect(testUser).NotTo(BeNil())
		Expect(testUser.ID).NotTo(BeZero())
	})

	AfterEach(func() {
		err := tokenClient.DeleteUser(ctx, testUser.ID)
		Expect(err).NotTo(HaveOccurred())
		deletedUser, err := tokenClient.GetUserByUsername(ctx, userInput.Username)
		Expect(err).To(HaveOccurred())
		Expect(deletedUser).To(BeNil())
	})

	It("can login", func() {
		config := urlfinder.ClientConfig{Lifecycle: lifecycle, Region: "na"}
		client := NewClient(config)
		err := client.Login(pixoUsername, pixoPassword)
		Expect(err).NotTo(HaveOccurred())
		Expect(client.IsAuthenticated()).To(BeTrue())
	})

	It("can get a user by username", func() {
		retrievedUser, err := tokenClient.GetUserByUsername(ctx, userInput.Username)
		Expect(err).NotTo(HaveOccurred())
		Expect(retrievedUser).NotTo(BeNil())
		Expect(retrievedUser.ID).To(Equal(testUser.ID))
		Expect(retrievedUser.Username).To(Equal(userInput.Username))
		Expect(retrievedUser.FirstName).To(Equal(userInput.FirstName))
		Expect(retrievedUser.LastName).To(Equal(userInput.LastName))
		Expect(retrievedUser.OrgID).To(Equal(userInput.OrgID))
	})

	It("can return an error if user does not exist when updating", func() {
		updatedUser, err := tokenClient.UpdateUser(ctx, userInput)
		Expect(err).To(HaveOccurred())
		Expect(updatedUser).To(BeNil())
		Expect(err.Error()).To(ContainSubstring("user id is required"))
	})

	It("can update a user", func() {
		userToUpdateInput := platform.User{
			FirstName: faker.FirstName(),
			LastName:  faker.LastName(),
			Username:  faker.Username(),
			Password:  faker.Password(),
			OrgID:     1,
		}
		userToUpdate, err := tokenClient.CreateUser(ctx, userToUpdateInput)
		Expect(err).NotTo(HaveOccurred())
		Expect(userToUpdate).NotTo(BeNil())
		Expect(userToUpdate.ID).NotTo(BeZero())
		userToUpdateInput.ID = userToUpdate.ID
		userToUpdateInput.Role = "superadmin"
		userToUpdateInput.FirstName = faker.FirstName()
		userToUpdateInput.LastName = faker.LastName()

		updatedUser, err := tokenClient.UpdateUser(ctx, userToUpdateInput)

		Expect(err).NotTo(HaveOccurred())
		Expect(updatedUser).NotTo(BeNil())
		Expect(updatedUser.ID).To(Equal(userToUpdateInput.ID))
		Expect(updatedUser.Username).To(Equal(userToUpdateInput.Username))
		Expect(updatedUser.FirstName).To(Equal(userToUpdateInput.FirstName))
		Expect(updatedUser.LastName).To(Equal(userToUpdateInput.LastName))
		Expect(updatedUser.Role).To(Equal(userToUpdateInput.Role))
	})

	Context("using an api key", func() {

		var (
			apiKey *platform.APIKey
		)

		BeforeEach(func() {
			var err error
			apiKey, err = tokenClient.CreateAPIKey(ctx, platform.APIKey{
				UserID: testUser.ID,
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(apiKey).NotTo(BeNil())
			Expect(apiKey.ID).NotTo(BeZero())
			Expect(apiKey.UserID).To(Equal(testUser.ID))
			Expect(apiKey.Key).NotTo(BeEmpty())
		})

		AfterEach(func() {
			err := tokenClient.DeleteAPIKey(ctx, apiKey.ID)
			Expect(err).NotTo(HaveOccurred())
		})

		It("can get a user with the api key", func() {
			config := urlfinder.ClientConfig{
				Lifecycle: lifecycle,
				Region:    "na",
				APIKey:    apiKey.Key,
			}
			client := NewClient(config)
			Expect(client.IsAuthenticated()).To(BeTrue())

			retrievedUser, err := client.GetUserByUsername(ctx, userInput.Username)
			Expect(err).NotTo(HaveOccurred())
			Expect(retrievedUser).NotTo(BeNil())
		})

		It("can get api keys", func() {
			apiKeys, err := tokenClient.GetAPIKeys(ctx, &APIKeyQueryParams{
				UserID: &testUser.ID,
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(apiKeys).NotTo(BeNil())
			Expect(len(apiKeys)).To(Equal(1))
			Expect(apiKeys[0].ID).To(Equal(apiKey.ID))
			Expect(apiKeys[0].UserID).To(Equal(testUser.ID))
		})

	})

})
