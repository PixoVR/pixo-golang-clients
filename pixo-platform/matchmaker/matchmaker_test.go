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
		m = matchmaker.NewMatchmaker("", os.Getenv("AUTH_TOKEN"))
	})

	It("can return an error message if the module ID is invalid", func() {
		req := matchmaker.MatchRequest{
			ModuleID:      0,
			OrgID:         1,
			ServerVersion: "1.00.00",
		}
		addr, err := m.Connect(req)

		Expect(err).To(HaveOccurred())
		Expect(addr).To(BeNil())
	})

	It("can return an error message if the org ID is invalid", func() {
		req := matchmaker.MatchRequest{
			ModuleID:      1,
			OrgID:         0,
			ServerVersion: "1.00.00",
		}
		addr, err := m.Connect(req)

		Expect(err).To(HaveOccurred())
		Expect(addr).To(BeNil())
	})

	It("can return an error message if the server version is invalid", func() {
		req := matchmaker.MatchRequest{
			ModuleID:      1,
			OrgID:         1,
			ServerVersion: "",
		}
		addr, err := m.Connect(req)

		Expect(err).To(HaveOccurred())
		Expect(addr).To(BeNil())
	})

	It("can retrieve a game server address using the matchmaker", func() {
		req := matchmaker.MatchRequest{
			ModuleID:      43,
			OrgID:         20,
			ServerVersion: "1.03.01",
		}
		addr, err := m.Connect(req)

		Expect(err).NotTo(HaveOccurred())
		Expect(addr).NotTo(BeNil())
		Expect(addr.IP).NotTo(BeEmpty())
		Expect(addr.Port).NotTo(BeZero())
	})

})
