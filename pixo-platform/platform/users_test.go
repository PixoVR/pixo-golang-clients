package platform_test

import (
	"context"
	"github.com/PixoVR/pixo-golang-clients/pixo-platform/platform"
	. "github.com/PixoVR/pixo-golang-clients/pixo-platform/platform"
	"github.com/PixoVR/pixo-golang-clients/pixo-platform/urlfinder"
	"github.com/go-faker/faker/v4"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Users", func() {

	var (
		ctx             context.Context
		user            *platform.User
		newUserPassword = faker.Password()
	)

	BeforeEach(func() {
		ctx = context.Background()

		user = &platform.User{
			FirstName: faker.FirstName(),
			LastName:  faker.LastName(),
			Username:  faker.Username(),
			Email:     faker.Email(),
			Password:  newUserPassword,
			OrgID:     1,
		}
		Expect(tokenClient.CreateUser(ctx, user)).To(Succeed())
		Expect(user).NotTo(BeNil())
		Expect(user.ID).NotTo(BeZero())
		Expect(user.Username).To(Equal(user.Username))
		Expect(user.Email).To(Equal(user.Email))
		user.Password = newUserPassword
	})

	AfterEach(func() {
		err := tokenClient.DeleteUser(ctx, user.ID)
		Expect(err).NotTo(HaveOccurred())
		deletedUser, err := tokenClient.GetUserByUsername(ctx, user.Username)
		Expect(err).To(HaveOccurred())
		Expect(deletedUser).To(BeNil())
	})

	It("can login", func() {
		config := urlfinder.ClientConfig{Lifecycle: lifecycle, Region: "na"}
		client := NewClient(config)

		err := client.Login(user.Username, user.Password)

		Expect(err).NotTo(HaveOccurred())
		Expect(client.IsAuthenticated()).To(BeTrue())
		Expect(client.ActiveUserID()).To(Equal(user.ID))
		Expect(client.ActiveOrgID()).To(Equal(user.OrgID))
	})

	It("can get a user by username", func() {
		retrievedUser, err := tokenClient.GetUserByUsername(ctx, user.Username)
		Expect(err).NotTo(HaveOccurred())
		Expect(retrievedUser).NotTo(BeNil())
		Expect(retrievedUser.ID).To(Equal(user.ID))
		Expect(retrievedUser.Username).To(Equal(user.Username))
		Expect(retrievedUser.FirstName).To(Equal(user.FirstName))
		Expect(retrievedUser.LastName).To(Equal(user.LastName))
		Expect(retrievedUser.OrgID).To(Equal(user.OrgID))
	})

	It("can return an error if user does not exist when updating", func() {
		err := tokenClient.UpdateUser(ctx, &platform.User{
			ID:       0,
			Username: faker.Username(),
		})
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("user id is required"))
	})

	It("can update a user", func() {
		updatedUser := &platform.User{
			ID:        user.ID,
			FirstName: faker.FirstName(),
			LastName:  faker.LastName(),
			Username:  faker.Username(),
			Password:  faker.Password(),
			OrgID:     1,
			Role:      "superadmin",
		}

		err := tokenClient.UpdateUser(ctx, updatedUser)

		Expect(err).NotTo(HaveOccurred())
		Expect(updatedUser).NotTo(BeNil())
		Expect(updatedUser.ID).To(Equal(updatedUser.ID))
		Expect(updatedUser.OrgID).To(Equal(updatedUser.OrgID))
		Expect(updatedUser.Username).To(Equal(updatedUser.Username))
		Expect(updatedUser.FirstName).To(Equal(updatedUser.FirstName))
		Expect(updatedUser.LastName).To(Equal(updatedUser.LastName))
		Expect(updatedUser.Role).To(Equal(updatedUser.Role))
	})

	Context("using an api key", func() {

		var (
			apiKey *APIKey
		)

		BeforeEach(func() {
			var err error
			apiKey, err = tokenClient.CreateAPIKey(ctx, APIKey{
				UserID: user.ID,
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(apiKey).NotTo(BeNil())
			Expect(apiKey.ID).NotTo(BeZero())
			Expect(apiKey.UserID).To(Equal(user.ID))
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

			retrievedUser, err := client.GetUserByUsername(ctx, user.Username)
			Expect(err).NotTo(HaveOccurred())
			Expect(retrievedUser).NotTo(BeNil())
		})

		It("can get api keys", func() {
			apiKeys, err := tokenClient.GetAPIKeys(ctx, &APIKeyQueryParams{
				UserID: &user.ID,
			})

			Expect(err).NotTo(HaveOccurred())
			Expect(apiKeys).NotTo(BeNil())
			Expect(len(apiKeys)).To(Equal(1))
			Expect(apiKeys[0].ID).To(Equal(apiKey.ID))
			Expect(apiKeys[0].UserID).To(Equal(user.ID))
		})

	})

})
