package matchmaker_test

import (
	"fmt"
	"github.com/PixoVR/pixo-golang-clients/pixo-platform/matchmaker"
	"github.com/PixoVR/pixo-golang-clients/pixo-platform/urlfinder"
	config2 "github.com/PixoVR/pixo-golang-server-utilities/pixo-platform/config"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"os"
)

var _ = Describe("Matchmaker", func() {

	var (
		m matchmaker.Matchmaker
	)

	BeforeEach(func() {
		var err error
		config := urlfinder.ClientConfig{
			Lifecycle: config2.GetEnvOrReturn("TEST_PIXO_LIFECYCLE", "dev"),
			Region:    config2.GetEnvOrReturn("TEST_PIXO_REGION", "na"),
		}
		m, err = matchmaker.NewClientWithBasicAuth(os.Getenv("TEST_PIXO_USERNAME"), os.Getenv("TEST_PIXO_PASSWORD"), config)
		Expect(err).NotTo(HaveOccurred())
	})

	It("can get the base url for the matchmaker", func() {
		Expect(m.GetURL("ws")).To(Equal(fmt.Sprintf("wss://apex.%s.pixovr.com/matchmaking", config2.GetEnvOrReturn("TEST_PIXO_LIFECYCLE", "stage"))))
	})

	It("can dial a the matchmaking service and request a match", func() {
		req := matchmaker.MatchRequest{
			ModuleID:      43,
			ServerVersion: "1.03.02",
		}

		conn, _, err := m.DialMatchmaker()
		Expect(err).NotTo(HaveOccurred())
		Expect(conn).NotTo(BeNil())

		err = m.SendRequest(conn, req)
		Expect(err).NotTo(HaveOccurred())

		resp, err := m.ReadResponse(conn)
		Expect(err).NotTo(HaveOccurred())
		Expect(resp).NotTo(BeNil())
		Expect(resp.Message).To(Equal("Match found"))

		err = m.CloseMatchmakerConnection(conn)
		Expect(err).NotTo(HaveOccurred())
	})

	It("can return an error message if the module ID is invalid", func() {
		req := matchmaker.MatchRequest{
			ModuleID:      0,
			ServerVersion: "1.00.00",
		}

		addr, err := m.FindMatch(req)

		Expect(err).To(HaveOccurred())
		Expect(addr).To(BeNil())
	})

	It("can return an error message if the server version is invalid", func() {
		req := matchmaker.MatchRequest{
			ModuleID:      1,
			ServerVersion: "",
		}

		addr, err := m.FindMatch(req)

		Expect(err).To(HaveOccurred())
		Expect(addr).To(BeNil())
	})

	It("can find a match and talk to the gameserver", func() {
		req := matchmaker.MatchRequest{
			ModuleID:      43,
			ServerVersion: "1.03.02",
		}

		addr, err := m.FindMatch(req)

		Expect(err).NotTo(HaveOccurred())
		Expect(addr).NotTo(BeNil())
		Expect(addr.IP).NotTo(BeEmpty())
		Expect(addr.Port).NotTo(BeZero())

		Expect(m.DialGameserver(addr)).To(Succeed())
		Expect(m.SendMessageToGameserver([]byte("hello world"))).To(Succeed())

		response, err := m.ReadMessageFromGameserver()
		Expect(err).NotTo(HaveOccurred())
		Expect(response).NotTo(BeEmpty())
	})

})
