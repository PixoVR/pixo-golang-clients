package cmd_test

import (
	"fmt"
	"github.com/go-faker/faker/v4"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Users Create", func() {

	BeforeEach(func() {
		executor = NewTestExecutor()
	})

	AfterEach(func() {
		executor.Cleanup()
	})

	It("can create a user", func() {
		firstName := faker.FirstName()
		lastName := faker.LastName()
		username := faker.Username()
		orgID := "1"
		role := "developer"
		password := faker.Password() + "!"

		output, err := executor.RunCommand(
			"users",
			"create",
			"--first-name",
			firstName,
			"--last-name",
			lastName,
			"--username",
			username,
			"--password",
			password,
			"--org-id",
			orgID,
			"--role",
			role,
		)
		Expect(err).NotTo(HaveOccurred())
		Expect(executor.MockPlatformClient.NumCalledCreateUser).To(Equal(1))
		Expect(output).To(ContainSubstring(fmt.Sprintf("User created: %s", username)))
	})

})
