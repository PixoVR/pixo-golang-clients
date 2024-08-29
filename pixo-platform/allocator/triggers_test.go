package allocator_test

//var _ = Describe("Triggers", Ordered, func() {
//
//	var (
//		allocatorClient *Client
//		trigger         = platform.MultiplayerServerTrigger{
//			ModuleID: 1,
//			Module: &platform.Module{
//				ID: 1,
//				GitConfig: platform.GitConfig{
//					OrgName:  "PixoVR",
//					RepoName: "multiplayer-gameservers",
//				},
//			},
//			Revision:   "dev",
//			Dockerfile: "Dockerfile",
//			Context:    "simple-server",
//			Config:     "version.ini",
//		}
//	)
//
//	BeforeEach(func() {
//		lifecycle := config2.GetEnvOrReturn("TEST_PIXO_LIFECYCLE", "dev")
//		apiKey := os.Getenv("TEST_PIXO_API_KEY")
//
//		config := urlfinder.ClientConfig{Lifecycle: lifecycle, APIKey: apiKey}
//		allocatorClient = NewClient(config)
//		Expect(allocatorClient.IsAuthenticated()).To(BeTrue())
//	})
//
//	It("can register a multiplayer server trigger", func() {
//		res := allocatorClient.RegisterTrigger(trigger)
//		Expect(res.Error).NotTo(HaveOccurred())
//		Expect(res.HTTPResponse.StatusCode()).To(Equal(http.StatusCreated))
//	})
//
//	It("can update a multiplayer server trigger", func() {
//		res := allocatorClient.UpdateTrigger(trigger)
//		Expect(res.Error).NotTo(HaveOccurred())
//		Expect(res.HTTPResponse.StatusCode()).To(Equal(http.StatusOK))
//	})
//
//	It("can delete a multiplayer server trigger", func() {
//		res := allocatorClient.DeleteTrigger(trigger.ID)
//		Expect(res.Error).NotTo(HaveOccurred())
//		Expect(res.HTTPResponse.StatusCode()).To(Equal(http.StatusNoContent))
//	})
//
//})
