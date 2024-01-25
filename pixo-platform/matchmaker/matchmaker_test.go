package matchmaker_test

import (
	"github.com/PixoVR/pixo-golang-clients/pixo-platform/matchmaker"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"os"
)

var _ = Describe("Multiplayer", func() {

	var (
		m matchmaker.Matchmaker
	)

	BeforeEach(func() {
		var err error
		m, err = matchmaker.NewMatchmakerWithBasicAuth(os.Getenv("PIXO_USERNAME"), os.Getenv("PIXO_PASSWORD"), "dev", "na")
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

	It("can retrieve a game server address using the matchmaker and send a message to it", func() {
		req := matchmaker.MatchRequest{
			ModuleID:      271,
			ServerVersion: "2.00.01",
		}

		addr, err := m.FindMatch(req)

		Expect(err).NotTo(HaveOccurred())
		Expect(addr).NotTo(BeNil())
		Expect(addr.IP).NotTo(BeEmpty())
		Expect(addr.Port).NotTo(BeZero())

		err = m.DialGameserver(addr)
		Expect(err).NotTo(HaveOccurred())

		response, err := m.SendAndReceiveMessage([]byte("hello world"))
		Expect(err).NotTo(HaveOccurred())
		Expect(response).NotTo(BeEmpty())
	})

})
