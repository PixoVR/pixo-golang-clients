package matchmaker_test

import (
	"github.com/PixoVR/pixo-golang-clients/pixo-platform/matchmaker"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"os"
)

var _ = Describe("Multiplayer", func() {

	var (
		matchURL = "wss://match.apex.pixovr.com/matchmaking/matchmake"
		m        matchmaker.Matchmaker
	)

	BeforeEach(func() {
		m = matchmaker.NewMatchmaker(matchURL, os.Getenv("AUTH_TOKEN"))
	})

	It("can return an error message if the module ID is invalid", func() {
		addr, err := m.Connect(99999999, 20)

		Expect(err).To(HaveOccurred())
		Expect(addr).To(BeNil())
	})

	It("can retrieve a game server address using the matchmaker", func() {
		addr, err := m.Connect(43, 20)

		Expect(err).NotTo(HaveOccurred())
		Expect(addr).NotTo(BeNil())
		Expect(addr.IP).NotTo(BeEmpty())
		Expect(addr.Port).NotTo(BeZero())
	})

})
