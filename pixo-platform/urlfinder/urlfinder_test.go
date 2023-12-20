package urlfinder_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/PixoVR/pixo-golang-clients/pixo-platform/urlfinder"
)

var _ = Describe("Urlfinder", func() {

	Context("domaon", func() {

		It("can find the default domain", func() {
			config := urlfinder.ServiceConfig{}
			domain := config.GetBaseDomain()
			Expect(domain).To(Equal("pixovr.com"))
		})

		It("can use localhost", func() {
			config := urlfinder.ServiceConfig{Lifecycle: "local"}
			domain := config.GetBaseDomain()
			Expect(domain).To(Equal("localhost"))
		})

		It("can find the dev domain", func() {
			config := urlfinder.ServiceConfig{Lifecycle: "dev"}
			domain := config.GetBaseDomain()
			Expect(domain).To(Equal("dev.pixovr.com"))
		})

	})

	Context("platform api", func() {

		It("can find the default url for the pixo platform api locally", func() {
			config := urlfinder.ServiceConfig{
				Lifecycle: "local",
				Port:      3001,
			}
			url := config.FormatURL()
			Expect(url).To(Equal("http://localhost:3001"))
		})

		It("can find the default url for the pixo platform api and ignore the port since its not local", func() {
			config := urlfinder.ServiceConfig{Port: 8000}
			url := config.FormatURL()
			Expect(url).To(Equal("https://primary.apex.pixovr.com"))
		})

		It("can find the url for the NA dev platform API", func() {
			config := urlfinder.ServiceConfig{Lifecycle: "dev"}
			url := config.FormatURL()
			Expect(url).To(Equal("https://primary.apex.dev.pixovr.com"))
		})

		It("can find the url for the saudi dev platform API", func() {
			config := urlfinder.ServiceConfig{
				Region:    "saudi",
				Lifecycle: "dev",
			}
			url := config.FormatURL()
			Expect(url).To(Equal("https://saudi.primary.apex.dev.pixovr.com"))
		})

	})

	Context("matchmaking", func() {

		It("can find the local matchmaking url", func() {
			config := urlfinder.ServiceConfig{Service: "match", Lifecycle: "local", Port: 8080}
			url := config.FormatURL()
			Expect(url).To(Equal("ws://localhost:8080"))
		})

		It("can find the default matchmaking url", func() {
			config := urlfinder.ServiceConfig{Service: "match"}
			url := config.FormatURL()
			Expect(url).To(Equal("wss://match.apex.pixovr.com"))
		})

		It("can find the saudi matchmaking url", func() {
			config := urlfinder.ServiceConfig{
				Service:   "match",
				Region:    "saudi",
				Lifecycle: "stage",
			}
			url := config.FormatURL()
			Expect(url).To(Equal("wss://saudi.match.apex.stage.pixovr.com"))
		})

		It("can find the url for the allocator service", func() {
			config := urlfinder.ServiceConfig{
				Service:   "allocator",
				Tenant:    "multiplayer",
				Lifecycle: "dev",
			}
			url := config.FormatURL()
			Expect(url).To(Equal("https://multi-central1.allocator.multiplayer.dev.pixovr.com"))
		})

		It("can find the url for the saudi allocator service", func() {
			config := urlfinder.ServiceConfig{
				Service:   "allocator",
				Tenant:    "multiplayer",
				Region:    "saudi",
				Lifecycle: "dev",
			}
			url := config.FormatURL()
			Expect(url).To(Equal("https://multi-saudi.allocator.multiplayer.dev.pixovr.com"))
		})

	})

})
