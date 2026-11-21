package allocator

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"

	abstractClient "github.com/PixoVR/pixo-golang-clients/pixo-platform/abstract"
	"github.com/PixoVR/pixo-golang-clients/pixo-platform/urlfinder"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("GetGameServers", func() {

	var (
		status      int
		body        string
		lastRequest *http.Request
		client      *Client
	)

	BeforeEach(func() {
		status = http.StatusOK
		body = `{"count":0,"results":[]}`
		lastRequest = nil

		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			lastRequest = r
			w.WriteHeader(status)
			_, _ = w.Write([]byte(body))
		})
		server := httptest.NewServer(handler)
		DeferCleanup(server.Close)

		serverURL, err := url.Parse(server.URL)
		Expect(err).NotTo(HaveOccurred())
		port, err := strconv.Atoi(serverURL.Port())
		Expect(err).NotTo(HaveOccurred())

		serviceConfig := urlfinder.ServiceConfig{Service: "allocator", Lifecycle: "local", Port: port}
		abstractConfig := abstractClient.Config{ServiceConfig: serviceConfig, Token: "secret-key"}
		client = &Client{ServiceClient: *abstractClient.NewClient(abstractConfig)}
	})

	It("requests the game servers with the token as a bearer token", func() {
		_, err := client.GetGameServers(context.Background())

		Expect(err).NotTo(HaveOccurred())
		Expect(lastRequest.Method).To(Equal(http.MethodGet))
		Expect(lastRequest.URL.Path).To(Equal("/allocator/gameservers"))
		Expect(lastRequest.URL.RawQuery).To(BeEmpty())
		Expect(lastRequest.Header.Get("Authorization")).To(Equal("Bearer secret-key"))
	})

	It("returns the game servers in the response", func() {
		body = `{"count":2,"results":[
			{"resource_name":"gs-1","ipaddress":"10.0.0.1","port":"7777","org_id":5,"module_id":7,"state":"Allocated","server_version":"1.0.0","num_playing":2,"players":["a","b"],"capacity":8,"created_at":"2026-11-21T00:00:00Z"},
			{"resource_name":"gs-2","org_id":-1,"module_id":-1,"state":"Ready"}
		]}`

		gameServers, err := client.GetGameServers(context.Background())

		Expect(err).NotTo(HaveOccurred())
		Expect(gameServers).To(Equal([]GameServer{
			{Name: "gs-1", IP: "10.0.0.1", Port: "7777", OrgID: 5, ModuleID: 7, State: "Allocated", ServerVersion: "1.0.0", NumPlaying: 2, Players: []string{"a", "b"}, Capacity: 8, CreatedAt: "2026-11-21T00:00:00Z"},
			{Name: "gs-2", OrgID: -1, ModuleID: -1, State: "Ready"},
		}))
	})

	It("returns an empty list when there are no game servers", func() {
		gameServers, err := client.GetGameServers(context.Background())

		Expect(err).NotTo(HaveOccurred())
		Expect(gameServers).To(BeEmpty())
	})

	It("returns an error when the allocator does not respond with 200", func() {
		status = http.StatusUnauthorized
		body = `{"error":"unauthorized"}`

		gameServers, err := client.GetGameServers(context.Background())

		Expect(err).To(MatchError(ContainSubstring("401")))
		Expect(gameServers).To(BeNil())
	})

	It("returns an error when the response is not valid JSON", func() {
		body = `not json`

		gameServers, err := client.GetGameServers(context.Background())

		Expect(err).To(HaveOccurred())
		Expect(gameServers).To(BeNil())
	})
})
