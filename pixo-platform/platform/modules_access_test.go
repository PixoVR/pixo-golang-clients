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
		Expect(err).To(MatchError("record not found"))
	})

	It("can return the modules a user has access to", func() {
		modules, err := tokenClient.GetModulesForUser(ctx, user.ID)
		Expect(err).NotTo(HaveOccurred())
		Expect(len(modules)).To(BeNumerically(">", 0))
	})
})
