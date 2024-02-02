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
	"net/http"
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

	return &Tester{
		client:      config.MatchmakingClient,
		connections: config.Connections,
		duration:    config.Duration,
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

	fmt.Printf(emoji.Sprintf(":check_mark_button: Connection %d: %s\n", id, successColor.Sprintf("established")))

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
		log.Debug().Err(err).Bytes("response", messageBytes).Msg("error reading message")
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
	fmt.Println()
	fmt.Println(emoji.Sprintf(":rocket: Starting load test with %d connections to %s...\n", lt.connections, cyanColor.Sprint(lt.client.GetURL())))

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

	fmt.Printf(emoji.Sprintf(":fountain_pen: Connection %d: %s: %s\n", id, successColor.Sprint("received match"), statColor.Sprint(string(msg))))
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
	fmt.Printf(emoji.Sprintf(":exclamation: Connection %d: %s - %s\n", id, errorColor.Sprint(msg), err))
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

	headerColor.Println("\nMatchmaking Load Test Summary")
	fmt.Println("================================")

	fmt.Printf("Max Test Duration:       %s\n", lt.duration)
	fmt.Printf("Actual Test Duration:    %s\n", lt.end.Sub(lt.start).Round(50*time.Millisecond))
	fmt.Printf("Connections:             %d\n", lt.connections)
	fmt.Printf("\nTotal Messages Sent:     %d\n", lt.messagesSent)
	statColor.Printf("Total Messages Received: %d\n", lt.messagesReceived)
	errorColor.Printf("Connection Errors:       %d\n", len(lt.connectionErrors))
	errorColor.Printf("Matching Errors:         %d\n", len(lt.matchingErrors))
	successColor.Printf("Matches Received:        %d\n", lt.successMessagesReceived)

	fmt.Println()
	lineColor.Println("┌─────────────┬────────────┐")
	headerColor.Println("│ Stat        │ Value      │")
	lineColor.Println("├─────────────┼────────────┤")
	statColor.Printf("│ Avg Latency │ %.2f s    │\n", avgLatency)
	statColor.Printf("│ Max Latency │ %.2f s    │\n", float64(lt.maxLatency)/float64(time.Second))
	statColor.Printf("│ Msgs/Sec    │ %.2f       │\n", messagesPerSecond)
	lineColor.Println("└─────────────┴────────────┘")
	fmt.Println()
}
