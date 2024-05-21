package cmd_test

import (
	"fmt"
	"github.com/go-faker/faker/v4"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Users", func() {

	var (
		executor *TestExecutor
	)

	BeforeEach(func() {
		executor = NewTestExecutor()
	})

	AfterEach(func() {
		executor.Cleanup()
	})

	It("can create a user and login", func() {
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
		Expect(executor.MockPlatformClient.CalledCreateUser).To(BeTrue())
		Expect(output).To(ContainSubstring(fmt.Sprintf("User created: %s", username)))

		executor.ExpectLoginToSucceed(username, password)

		output, err = executor.RunCommand("config")
		Expect(output).To(ContainSubstring(fmt.Sprintf("Username: %s", username)))
		Expect(err).NotTo(HaveOccurred())
	})

})
