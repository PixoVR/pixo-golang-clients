package load

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/PixoVR/pixo-golang-clients/pixo-platform/matchmaker"
	"github.com/fatih/color"
	"github.com/gorilla/websocket"
	"github.com/kyokomi/emoji"
	"github.com/rs/zerolog/log"
	"io"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
)

var (
	headerColor  = color.New(color.FgHiCyan, color.Bold)
	cyanColor    = color.New(color.FgHiCyan)
	successColor = color.New(color.FgHiGreen)
	errorColor   = color.New(color.FgHiRed)
	statColor    = color.New(color.FgHiYellow)
	lineColor    = color.New(color.FgHiBlue)
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
	mu                      sync.Mutex
	reader                  io.Reader
	writer                  io.Writer
}

// NewLoadTester creates a new instance of Tester with the specified configuration.
func NewLoadTester(config Config) *Tester {
	if config.MatchmakingClient == nil {
		log.Fatal().Msg("matchmaking client is required")
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
		client:      config.MatchmakingClient,
		connections: config.Connections,
		duration:    config.Duration,
		writer:      config.Writer,
		reader:      config.Reader,
	}
}

// performRequest establishes a single WebSocket connection and requests a match
func (lt *Tester) performRequest(wg *sync.WaitGroup, id int) {
	defer wg.Done()

	httpHeader := http.Header{}
	httpHeader.Add("Authorization", "Bearer "+lt.client.GetToken())

	conn, _, err := websocket.DefaultDialer.Dial(lt.client.GetURL(), httpHeader)
	if err != nil {
		lt.recordConnectionError(id, "failed to connect", err)
		lt.recordFailure()
		return
	}
	defer conn.Close()
	lt.recordSuccess()

	if err = conn.SetReadDeadline(time.Now().Add(30 * time.Second)); err != nil {
		return
	}

	lt.print(emoji.Sprintf(":check_mark_button: Connection %d: %s\n", id, successColor.Sprintf("established")))

	reqBytes, err := json.Marshal(lt.request)
	if err != nil {
		lt.recordConnectionError(id, "failed to marshal request", err)
		return
	}

	start := time.Now()
	if err = conn.WriteMessage(websocket.TextMessage, reqBytes); err != nil {
		lt.recordConnectionError(id, "failed to send message", err)
		return
	}
	lt.recordSentMessage()

	_, messageBytes, err := conn.ReadMessage()
	lt.recordLatency(time.Since(start))
	lt.recordReceivedMessage()
	if err != nil {
		lt.recordConnectionError(id, "error reading message", err)
		return
	}

	if !strings.Contains(string(messageBytes), "IPAddress") || !strings.Contains(string(messageBytes), "Port") {
		lt.recordMatchingError(id, "matchmaking error", errors.New(string(messageBytes)))
		return
	}

	lt.recordSuccessMessageReceived(id, messageBytes)
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

// recordSuccessMessageReceived increments the count of successful messages received.
func (lt *Tester) recordSuccessMessageReceived(id int, msg []byte) {
	lt.mu.Lock()
	defer lt.mu.Unlock()

	lt.print(emoji.Sprintf(":fountain_pen: Connection %d: %s: %s\n", id, successColor.Sprint("received match"), statColor.Sprint(string(msg))))
	lt.successMessagesReceived++
}

// recordSuccess increments the count of successful connections.
func (lt *Tester) recordSuccess() {
	lt.mu.Lock()
	defer lt.mu.Unlock()
	lt.connectionSuccesses++
}

// recordFailure increments the count of failed connections.
func (lt *Tester) recordFailure() {
	lt.mu.Lock()
	defer lt.mu.Unlock()
	lt.connectionFailures++
}

func (lt *Tester) logError(id int, msg string, err error) {
	lt.print(emoji.Sprintf(":exclamation: Connection %d: %s - %s\n", id, errorColor.Sprint(msg), err))
}

// recordConnectionError adds an error to the list of encountered connectionErrors.
func (lt *Tester) recordConnectionError(id int, msg string, err error) {
	lt.mu.Lock()
	defer lt.mu.Unlock()

	lt.logError(id, msg, err)
	lt.connectionErrors = append(lt.connectionErrors, fmt.Sprintf("%s - %s", msg, err.Error()))
}

// recordMatchingError adds an error to the list of encountered matchingErrors.
func (lt *Tester) recordMatchingError(id int, msg string, err error) {
	lt.mu.Lock()
	defer lt.mu.Unlock()

	lt.logError(id, msg, err)
	lt.matchingErrors = append(lt.matchingErrors, fmt.Sprintf("%s - %s", msg, err.Error()))
}

// recordSentMessage increments the count of sent messages.
func (lt *Tester) recordSentMessage() {
	lt.mu.Lock()
	defer lt.mu.Unlock()

	lt.messagesSent++
}

// recordReceivedMessage increments the count of received messages.
func (lt *Tester) recordReceivedMessage() {
	lt.mu.Lock()
	defer lt.mu.Unlock()

	lt.messagesReceived++
}

// recordLatency tracks the latency of each message and updates max latency if necessary.
func (lt *Tester) recordLatency(latency time.Duration) {
	lt.mu.Lock()
	defer lt.mu.Unlock()

	lt.totalLatency += latency
	if latency > lt.maxLatency {
		lt.maxLatency = latency
	}

	lt.numLatencies++
}

// displayStats prints the collected statistics to the console.
func (lt *Tester) displayStats() {
	totalDuration := lt.end.Sub(lt.start).Seconds()
	var avgLatency float64
	if lt.numLatencies > 0 {
		avgLatency = float64(lt.totalLatency) / float64(lt.numLatencies) / float64(time.Second)
	}
	var messagesPerSecond float64
	if totalDuration > 0 {
		messagesPerSecond = float64(lt.messagesSent) / totalDuration
	}

	lt.print(headerColor.Sprint("\nMatchmaking Load Test Summary"))
	lt.print("================================")

	lt.print("Max Test Duration:       %s\n", lt.duration)
	lt.print("Actual Test Duration:    %s\n", lt.end.Sub(lt.start).Round(50*time.Millisecond))
	lt.println("Connections:             %d\n", lt.connections)
	lt.print("Total Messages Sent:     %d\n", lt.messagesSent)

	lt.print(statColor.Sprint("Total Messages Received: %d\n", lt.messagesReceived))
	lt.print(errorColor.Sprint("Connection Errors:       %d\n", len(lt.connectionErrors)))
	lt.print(errorColor.Sprint("Matching Errors:         %d\n", len(lt.matchingErrors)))
	lt.print(successColor.Sprint("Matches Received:        %d\n", lt.successMessagesReceived))

	lt.println()
	lt.print(lineColor.Sprint("┌─────────────┬────────────┐"))
	lt.print(headerColor.Sprint("│ Stat        │ Value      │"))
	lt.print(lineColor.Sprint("├─────────────┼────────────┤"))
	lt.print(statColor.Sprint("│ Avg Latency │ %.2f s    │\n", avgLatency))
	lt.print(statColor.Sprint("│ Max Latency │ %.2f s    │\n", float64(lt.maxLatency)/float64(time.Second)))
	lt.print(statColor.Sprint("│ Msgs/Sec    │ %.2f       │\n", messagesPerSecond))
	lt.print(lineColor.Sprint("└─────────────┴────────────┘"))
	lt.println()
}

func (lt *Tester) print(msgs ...interface{}) {
	for _, msg := range msgs {
		_, _ = lt.writer.Write([]byte(fmt.Sprint(msg)))
	}
}

func (lt *Tester) println(msgs ...interface{}) {
	lt.print(msgs...)
	lt.print("\n")
}
