package load

import (
	"errors"
	"github.com/PixoVR/pixo-golang-clients/pixo-platform/matchmaker"
	"github.com/kyokomi/emoji"
	"io"
	"os"
	"sync"
	"time"
)

// Config contains the configuration for a load test.
type Config struct {
	Reader            io.Reader
	Writer            io.Writer
	MatchmakingClient matchmaker.Matchmaker
	Request           matchmaker.MatchRequest
	Connections       int
	Duration          time.Duration
}

// Tester configures and runs WebSocket load tests.
type Tester struct {
	client                  matchmaker.Matchmaker
	request                 matchmaker.MatchRequest
	connections             int
	duration                time.Duration
	connectionSuccesses     int
	connectionFailures      int
	messagesSent            int
	messagesReceived        int
	successMessagesReceived int
	numLatencies            int
	totalLatency            time.Duration
	maxLatency              time.Duration
	start                   time.Time
	end                     time.Time
	connectionErrors        []string
	matchingErrors          []string
	gameserversCount        map[string]int
	mu                      sync.Mutex
	reader                  io.Reader
	writer                  io.Writer
}

// NewLoadTester creates a new instance of Tester with the specified configuration.
func NewLoadTester(config Config) (*Tester, error) {
	if config.MatchmakingClient == nil {
		return nil, errors.New("matchmaking client is required")
	}

	if !config.Request.IsValid() {
		return nil, errors.New("match request is invalid")
	}

	if config.Connections <= 0 {
		config.Connections = 1
	}

	if config.Duration <= 0 {
		config.Duration = 600 * time.Second
	}

	if config.Reader == nil {
		config.Reader = os.Stdin
	}

	if config.Writer == nil {
		config.Writer = os.Stdout
	}

	return &Tester{
		request:     config.Request,
		client:      config.MatchmakingClient,
		connections: config.Connections,
		duration:    config.Duration,
		writer:      config.Writer,
		reader:      config.Reader,
	}, nil
}

// Run starts the load testing process.
func (lt *Tester) Run() {
	lt.println()
	lt.println(emoji.Sprintf(":rocket: Starting load test with %d connections to %s...\n", lt.connections, cyanColor.Sprint(lt.client.GetURL())))

	lt.start = time.Now()

	var wg sync.WaitGroup

	for i := 0; i < lt.connections; i++ {
		wg.Add(1)
		go lt.performRequest(&wg, i+1)
	}

	wg.Wait()
	lt.end = time.Now()
	lt.displayStats()
}
