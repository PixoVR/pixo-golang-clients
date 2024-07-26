package headset_test

import (
	. "github.com/PixoVR/pixo-golang-clients/pixo-platform/headset"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Headset Client", func() {

	var (
		headsetClient Client
	)

	BeforeEach(func() {
		var err error
		headsetClient, err = NewClientWithBasicAuth(username, password, clientConfig)
		Expect(err).NotTo(HaveOccurred())
		Expect(headsetClient).NotTo(BeNil())
		Expect(headsetClient.IsAuthenticated()).To(BeTrue())
	})

	It("can login", func() {
		anonymousClient := NewClient(clientConfig)
		Expect(anonymousClient.IsAuthenticated()).To(BeFalse())
		Expect(anonymousClient.Login(username, password)).NotTo(HaveOccurred())
		Expect(anonymousClient.IsAuthenticated()).To(BeTrue())
	})

	//It("should throw an error if the session module exist", func() {
	//	session := &platform.Session{
	//		ModuleID: moduleID,
	//	}
	//	err := headsetClient.StartSession(9999999999)
	//	Expect(err).To(HaveOccurred())
	//	Expect(err.Error()).To(ContainSubstring("invalid session"))
	//	Expect(session).NotTo(BeNil())
	//})

})
