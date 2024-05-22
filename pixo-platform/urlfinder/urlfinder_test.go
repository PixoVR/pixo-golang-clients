package urlfinder_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/PixoVR/pixo-golang-clients/pixo-platform/urlfinder"
)

var _ = Describe("Urlfinder", func() {

	Context("default domain", func() {

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

		It("can use internal k8s dns", func() {
			config := urlfinder.ServiceConfig{
				Lifecycle:   "dev",
				InternalDNS: true,
				Namespace:   "dev-apex",
				ServiceName: "primary-api",
			}
			domain := config.FormatURL()
			Expect(domain).To(Equal("http://dev-apex-primary-api.dev-apex.svc"))
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
			}
			url := config.FormatURL()
			Expect(url).To(Equal("http://localhost:8000/v2"))
		})

		It("can find the default url for the pixo platform api and ignore the port since its not local", func() {
			config := urlfinder.ServiceConfig{Port: 8000}
			url := config.FormatURL()
			Expect(url).To(Equal("https://apex.pixovr.com/v2"))
		})

		It("can find the url for the NA dev platform API", func() {
			config := urlfinder.ServiceConfig{Lifecycle: "dev"}
			url := config.FormatURL()
			Expect(url).To(Equal("https://apex.dev.pixovr.com/v2"))
		})

		It("can find the url for the saudi dev platform API", func() {
			config := urlfinder.ServiceConfig{
				Region:    "saudi",
				Lifecycle: "dev",
			}
			url := config.FormatURL()
			Expect(url).To(Equal("https://saudi.apex.dev.pixovr.com/v2"))
		})

		It("can find the url for the legacy dev platform API", func() {
			config := urlfinder.ServiceConfig{
				Lifecycle: "dev",
				Service:   "api",
			}
			url := config.FormatURL()
			Expect(url).To(Equal("https://api.apex.dev.pixovr.com"))
		})

		It("can find the local old api url", func() {
			config := urlfinder.ServiceConfig{Service: "api", Lifecycle: "local", Port: 8003}
			url := config.FormatURL()
			Expect(url).To(Equal("http://localhost:8003"))
		})

		It("can find the url for the legacy saudi prod platform API", func() {
			config := urlfinder.ServiceConfig{
				Region:  "saudi",
				Service: "api",
			}
			url := config.FormatURL()
			Expect(url).To(Equal("https://saudi.api.apex.pixovr.com"))
		})

		It("can find the url for the saudi prod platform API", func() {
			config := urlfinder.ServiceConfig{
				Region: "saudi",
			}
			url := config.FormatURL()
			Expect(url).To(Equal("https://saudi.apex.pixovr.com/v2"))
		})

	})

	Context("headset api", func() {

		It("can find the local headset api url", func() {
			config := urlfinder.ServiceConfig{Service: "modules", Lifecycle: "local", Port: 8003}
			url := config.FormatURL()
			Expect(url).To(Equal("http://localhost:8003"))
		})

		It("can find the default headset api url", func() {
			config := urlfinder.ServiceConfig{Service: "modules"}
			url := config.FormatURL()
			Expect(url).To(Equal("https://modules.apex.pixovr.com"))
		})

	})

	Context("matchmaking", func() {

		It("can find the local matchmaking url", func() {
			config := urlfinder.ServiceConfig{
				Service:   "matchmaking",
				Lifecycle: "local",
				Port:      8080,
			}

			url := config.FormatURL()

			Expect(url).To(Equal("http://localhost:8080/matchmaking"))
		})

		It("can find the local matchmaking websocket url", func() {
			config := urlfinder.ServiceConfig{Service: "matchmaking", Lifecycle: "local", Port: 8080}
			url := config.FormatURL("ws")
			Expect(url).To(Equal("ws://localhost:8080/matchmaking"))
		})

		It("can generate the default matchmaking url", func() {
			config := urlfinder.ServiceConfig{Service: "matchmaking"}
			url := config.FormatURL()
			Expect(url).To(Equal("https://apex.pixovr.com/matchmaking"))
		})

		It("can generate the saudi matchmaking https url", func() {
			config := urlfinder.ServiceConfig{
				Service: "matchmaking",
				Region:  "saudi",
			}

			url := config.FormatURL()

			Expect(url).To(Equal("https://saudi.apex.pixovr.com/matchmaking"))
		})

		It("can find the stage saudi matchmaking url", func() {
			config := urlfinder.ServiceConfig{
				Service:   "matchmaking",
				Region:    "saudi",
				Lifecycle: "stage",
			}

			url := config.FormatURL()

			Expect(url).To(Equal("https://saudi.apex.stage.pixovr.com/matchmaking"))
		})

		It("can generate the saudi matchmaking websocket url", func() {
			config := urlfinder.ServiceConfig{
				Service: "matchmaking",
				Region:  "saudi",
			}

			url := config.FormatURL("ws")

			Expect(url).To(Equal("wss://saudi.apex.pixovr.com/matchmaking"))
		})

		It("can find the saudi matchmaking websocket url when formatting", func() {
			config := urlfinder.ServiceConfig{
				Service:   "matchmaking",
				Region:    "saudi",
				Lifecycle: "stage",
			}

			url := config.FormatURL("ws")

			Expect(url).To(Equal("wss://saudi.apex.stage.pixovr.com/matchmaking"))
		})

		It("can find the url for the allocator service", func() {
			config := urlfinder.ServiceConfig{
				Service:   "allocator",
				Tenant:    "multiplayer",
				Lifecycle: "dev",
			}

			url := config.FormatURL()

			Expect(url).To(Equal("https://multi-central1.multiplayer.dev.pixovr.com/allocator"))
		})

		It("can find the url for the saudi allocator service", func() {
			config := urlfinder.ServiceConfig{
				Service:   "allocator",
				Tenant:    "multiplayer",
				Region:    "saudi",
				Lifecycle: "dev",
			}

			url := config.FormatURL()

			Expect(url).To(Equal("https://multi-saudi.multiplayer.dev.pixovr.com/allocator"))
		})

	})

})
