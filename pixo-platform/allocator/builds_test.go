package allocator_test

import (
	. "github.com/onsi/ginkgo/v2"
)

var _ = Describe("Builds", func() {

	//var (
	//	allocatorClient *multiplayerAllocator.Client
	//)
	//
	//BeforeEach(func() {
	//	var err error
	//	config := urlfinder.ClientConfig{
	//		Lifecycle: "dev",
	//		Token:     os.Getenv("SECRET_KEY"),
	//	}
	//	allocatorClient = multiplayerAllocator.NewClient(config)
	//	Expect(err).NotTo(HaveOccurred())
	//	Expect(allocatorClient.IsAuthenticated()).To(BeTrue())
	//})
	//
	//It("can get the build workflows and their logs", func() {
	//	workflows, err := allocatorClient.GetBuildWorkflows()
	//	Expect(err).NotTo(HaveOccurred())
	//	Expect(workflows).NotTo(BeNil())
	//	Expect(len(workflows)).To(BeNumerically(">", 0))
	//
	//	logsCh, err := allocatorClient.GetBuildWorkflowLogs(workflows[0].Name)
	//	Expect(err).NotTo(HaveOccurred())
	//	Expect(logsCh).NotTo(BeNil())
	//	log := <-logsCh
	//	Expect(log).NotTo(BeNil())
	//	Expect(log.Step).NotTo(BeEmpty())
	//	Expect(log.Lines).NotTo(BeEmpty())
	//})

})
