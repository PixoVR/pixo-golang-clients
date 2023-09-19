package auth_test

import (
	"github.com/PixoVR/pixo-golang-clients/pixo-platform/middleware/auth"
	platform "github.com/PixoVR/pixo-golang-clients/pixo-platform/primary-api"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Auth", func() {
	It("can initialize a custom context with a findUserByID function", func() {
		customCtx := auth.CustomContext{
			FindUserByID: func(id int) (*interface{}, error) {
				user := &platform.User{
					ID: id,
				}
				userVal := interface{}(user)
				return &userVal, nil
			},
		}

		user, err := customCtx.FindUserByID(1)

		Expect(err).NotTo(HaveOccurred())
		Expect(user).NotTo(BeNil())
	})
})
