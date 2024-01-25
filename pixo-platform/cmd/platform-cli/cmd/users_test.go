package cmd_test

import (
	"fmt"
	"github.com/go-faker/faker/v4"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Users", func() {

	It("can create a user", func() {
		firstName := faker.FirstName()
		lastName := faker.LastName()
		username := faker.Username()
		orgID := "1"
		role := "developer"

		output, err := RunCommand(
			"users",
			"create",
			"--first-name",
			firstName,
			"--last-name",
			lastName,
			"--username",
			username,
			"--password",
			"SomePassword123!",
			"--org-id",
			orgID,
			"--role",
			role,
		)
		Expect(err).NotTo(HaveOccurred())
		Expect(output).To(ContainSubstring(fmt.Sprintf("Created user %s", username)))
	})

})
