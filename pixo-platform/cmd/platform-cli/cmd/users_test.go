package cmd_test

import (
	"fmt"
	"github.com/go-faker/faker/v4"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"io"
)

var _ = Describe("Users", func() {

	It("can create a user", func() {
		firstName := faker.FirstName()
		lastName := faker.LastName()
		username := faker.Username()
		orgID := "1"
		role := "developer"

		rootCmd, output := GetRootCmd()
		rootCmd.SetArgs([]string{
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
		})
		err := rootCmd.Execute()
		Expect(err).NotTo(HaveOccurred())

		out, err := io.ReadAll(output)
		Expect(err).NotTo(HaveOccurred())
		Expect(string(out)).To(ContainSubstring(fmt.Sprintf("created user %s", username)))
	})

})
