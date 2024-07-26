package cmd_test

import (
	"bytes"
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

	It("return an error if unable to get orgs", func() {
		executor.MockPlatformClient.GetOrgsError = fmt.Errorf("failed to get orgs")
		input := bytes.NewBufferString("\n")

		_, err := executor.RunCommandWithInput(
			input,
			"users",
			"create",
			"--first-name",
			"test",
			"--last-name",
			"test",
			"--user-email",
			"test",
			"--user-password",
			"test",
			"--role",
			"developer",
		)

		Expect(err).To(HaveOccurred())
		Expect(err).To(MatchError("failed to get orgs"))
		Expect(executor.MockPlatformClient.NumCalledGetOrgs).To(Equal(1))
		Expect(executor.MockPlatformClient.NumCalledCreateOrg).To(Equal(0))
	})

	It("return an error if unable to get roles", func() {
		executor.MockPlatformClient.GetRolesError = fmt.Errorf("failed to get roles")
		input := bytes.NewBufferString("\n")

		_, err := executor.RunCommandWithInput(
			input,
			"users",
			"create",
			"--first-name",
			"test",
			"--last-name",
			"test",
			"--user-email",
			"test",
			"--user-password",
			"test",
			"--org",
			"1",
		)

		Expect(err).To(HaveOccurred())
		Expect(err).To(MatchError("failed to get roles"))
		Expect(executor.MockPlatformClient.NumCalledGetRoles).To(Equal(1))
		Expect(executor.MockPlatformClient.NumCalledCreateOrg).To(Equal(0))
	})

	It("returns an error if the password is not provided", func() {
		input := bytes.NewBufferString("\n")

		_, err := executor.RunCommandWithInput(
			input,
			"users",
			"create",
			"--first-name",
			"test",
			"--last-name",
			"test",
			"--user-email",
			"test",
			"--user-username",
			"test",
			"--org",
			"1",
			"--role",
			"developer",
		)

		Expect(err).To(HaveOccurred())
		Expect(err).To(MatchError("password is required"))
		Expect(executor.MockPlatformClient.NumCalledGetOrgs).To(Equal(0))
	})

	It("returns an error if the create call fails", func() {
		executor.MockPlatformClient.CreateUserError = fmt.Errorf("failed to create user")
		input := bytes.NewBufferString("\n")

		_, err := executor.RunCommandWithInput(
			input,
			"users",
			"create",
			"--first-name",
			"test",
			"--last-name",
			"test",
			"--user-email",
			"test",
			"--user-username",
			"test",
			"--user-password",
			"test",
			"--org",
			"1",
			"--role",
			"developer",
		)

		Expect(err).To(HaveOccurred())
		Expect(err).To(MatchError("failed to create user"))
	})

	It("can create a user", func() {
		firstName := faker.FirstName()
		lastName := faker.LastName()
		email := faker.Email()
		username := faker.Username()
		orgID := "1"
		role := "developer"
		password := faker.Password()

		output, err := executor.RunCommand(
			"users",
			"create",
			"--first-name",
			firstName,
			"--last-name",
			lastName,
			"--user-username",
			username,
			"--user-email",
			email,
			"--user-password",
			password,
			"--org",
			orgID,
			"--role",
			role,
		)
		Expect(err).NotTo(HaveOccurred())
		Expect(executor.MockPlatformClient.NumCalledCreateUser).To(Equal(1))
		Expect(output).To(ContainSubstring(fmt.Sprintf("User created: %s - %s", email, role)))
	})

})
