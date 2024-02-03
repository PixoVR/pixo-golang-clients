package load

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/PixoVR/pixo-golang-clients/pixo-platform/matchmaker"
	"github.com/fatih/color"
	"github.com/kyokomi/emoji"
	"io"
	"net"
	"os"
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

// performRequest establishes a single WebSocket connection and requests a match
func (lt *Tester) performRequest(wg *sync.WaitGroup, id int) {
	defer wg.Done()

	message, err := json.Marshal(lt.request)
	if err != nil {
		return
	}

	if _, err = lt.client.DialWebsocket(matchmaker.MatchmakingEndpoint); err != nil {
		lt.recordConnectionError(id, "failed to connect", err)
		return
	}

	defer func() {
		_ = lt.client.CloseWebsocketConnection()
	}()

	lt.recordConnectionSuccess(id)

	start := time.Now()
	if err = lt.client.WriteToWebsocket(message); err != nil {
		lt.recordConnectionError(id, "failed to send message", err)
		return
	}

	lt.recordSentMessage()

	_, msg, err := lt.client.ReadFromWebsocket()
	lt.recordLatency(time.Since(start))
	lt.recordReceivedMessage()
	if err != nil {
		lt.recordConnectionError(id, "failed to read message", err)
		return
	}

	var matchResponse matchmaker.MatchResponse
	if err = json.Unmarshal(msg, &matchResponse); err != nil {
		lt.recordMatchingError(id, "failed to unmarshal response", err)
		return
	}

	if !matchResponse.IsValid() {
		lt.recordMatchingError(id, "received invalid match", errors.New(matchResponse.Message))
		return
	}

	lt.recordSuccessMessageReceived(id, matchResponse)
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
func (lt *Tester) recordSuccessMessageReceived(id int, response matchmaker.MatchResponse) {
	lt.mu.Lock()
	defer lt.mu.Unlock()

	addr := net.JoinHostPort(response.MatchDetails.IP, response.MatchDetails.Port)
	msg := fmt.Sprintf("gameserver -> %s", addr)
	lt.print(emoji.Sprintf(":checkered_flag:Connection %d: %s - %s\n", id, successColor.Sprint(response.Message), statColor.Sprint(msg)))
	lt.successMessagesReceived++
}

// recordConnectionSuccess increments the count of successful connections.
func (lt *Tester) recordConnectionSuccess(id int) {
	lt.mu.Lock()
	defer lt.mu.Unlock()

	lt.print(emoji.Sprintf(":check_mark_button: Connection %d: %s\n", id, successColor.Sprint("established")))
	lt.connectionSuccesses++
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
	lt.connectionFailures++
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

	lt.println(headerColor.Sprint("\nMatchmaking Load Test Summary"))
	lt.println("==============================")

	lt.printf("Max Test Duration:       %s", lt.duration)
	lt.printf("Actual Test Duration:    %s", lt.end.Sub(lt.start).Round(50*time.Millisecond))
	lt.printf("Connections:             %d", lt.connections)
	lt.printf("Total Messages Sent:     %d", lt.messagesSent)

	lt.println(statColor.Sprintf("\nTotal Messages Received: %d", lt.messagesReceived))
	lt.println(errorColor.Sprintf("Connection Errors:       %d", len(lt.connectionErrors)))
	lt.println(errorColor.Sprintf("Matching Errors:         %d", len(lt.matchingErrors)))
	lt.println(successColor.Sprintf("Matches Received:        %d", lt.successMessagesReceived))

	lt.println()
	lt.println(lineColor.Sprint("┌─────────────┬────────────┐"))
	lt.println(headerColor.Sprint("│ Stat        │ Value      │"))
	lt.println(lineColor.Sprint("├─────────────┼────────────┤"))
	lt.println(statColor.Sprintf("│ Avg Latency │ %.2f s    │", avgLatency))
	lt.println(statColor.Sprintf("│ Max Latency │ %.2f s    │", float64(lt.maxLatency)/float64(time.Second)))
	lt.println(statColor.Sprintf("│ Msgs/Sec    │ %.2f       │", messagesPerSecond))
	lt.println(lineColor.Sprint("└─────────────┴────────────┘"))
	lt.println()
}

func (lt *Tester) print(msgs ...interface{}) {
	for _, msg := range msgs {
		_, _ = lt.writer.Write([]byte(fmt.Sprint(msg)))
	}
}

func (lt *Tester) printf(format string, msgs ...interface{}) {
	lt.println(fmt.Sprintf(format, msgs...))
}

func (lt *Tester) println(msgs ...interface{}) {
	lt.print(msgs...)
	lt.print("\n")
}
