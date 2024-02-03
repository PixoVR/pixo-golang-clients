package load

import (
	"fmt"
	"time"
)

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
