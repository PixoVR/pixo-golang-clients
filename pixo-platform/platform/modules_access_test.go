package platform_test

import (
	"github.com/go-faker/faker/v4"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"context"

	"github.com/PixoVR/pixo-golang-clients/pixo-platform/platform"
)

var _ = Describe("ModulesAccess", func() {

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

		user.Password = newUserPassword
	})

	AfterEach(func() {
		err := tokenClient.DeleteUser(ctx, user.ID)
		Expect(err).NotTo(HaveOccurred())
		deletedUser, err := tokenClient.GetUserByUsername(ctx, user.Username)
		Expect(err).To(HaveOccurred())
		Expect(deletedUser).To(BeNil())
	})

	It("can fail to get the modules a user has access to if the user does not exist", func() {
		_, err := tokenClient.GetModulesForUser(ctx, 999999)
		Expect(err).To(HaveOccurred())
	})

	It("can return the modules a user has access to", func() {
		modules, err := tokenClient.GetModulesForUser(ctx, user.ID)
		Expect(err).NotTo(HaveOccurred())
		Expect(len(modules)).To(BeNumerically(">", 0))
	})

	It("can return the modules of several users in one request", func() {
		usersModules, err := tokenClient.GetModulesForUsers(ctx, []int{user.ID})
		Expect(err).NotTo(HaveOccurred())
		Expect(usersModules).To(HaveLen(1))
		Expect(usersModules[0].UserID).To(Equal(user.ID))
		Expect(len(usersModules[0].Modules)).To(BeNumerically(">", 0))
	})

	It("can fail to get the modules of several users if one does not exist", func() {
		_, err := tokenClient.GetModulesForUsers(ctx, []int{user.ID, 999999})
		Expect(err).To(HaveOccurred())
	})

	It("can return the users of an org that can access a module", func() {
		users, err := tokenClient.GetUsersWithModuleAccess(ctx, moduleID, orgID)

		Expect(err).NotTo(HaveOccurred())
		Expect(users).NotTo(BeEmpty())
		for _, moduleUser := range users {
			Expect(moduleUser.ID).NotTo(BeZero())
			Expect(moduleUser.OrgID).To(Equal(orgID))
		}
	})

	It("can fail to get the users that can access a module of an org it cannot read", func() {
		_, err := tokenClient.GetUsersWithModuleAccess(ctx, moduleID, 0)

		Expect(err).To(HaveOccurred())
	})
})
