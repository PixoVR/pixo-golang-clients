package load

import (
	"fmt"
	"github.com/PixoVR/pixo-golang-clients/pixo-platform/matchmaker"
	"github.com/kyokomi/emoji"
	"net"
	"time"
)

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
