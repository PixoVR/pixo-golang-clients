package fixture

import (
	"fmt"
	"strings"
	"time"
)

// DisplaySummary prints the collected statistics to the console.
func (t *Tester) DisplaySummary() {
	t.printErrors()

	totalDuration := t.end.Sub(t.start).Seconds()
	var avgLatency float64
	if t.NumLatencies > 0 {
		avgLatency = float64(t.TotalLatency) / float64(t.NumLatencies) / float64(time.Second)
	}
	var messagesPerSecond float64
	if totalDuration > 0 {
		messagesPerSecond = float64(t.NumDone) / totalDuration
	}

	t.Println(HeaderColor.Sprint("\nLoad Test Summary"))
	t.Println("===========================")

	t.PrintInt("Concurrent Workers", t.Config.Concurrent)
	t.PrintInt("Amount Requested", t.Config.Amount)
	t.PrintInt("Amount Completed", t.NumDone)
	t.PrintDuration("Max Test Duration", t.Config.MaxDuration)
	t.PrintDuration("Actual Test Duration", t.end.Sub(t.start))

	t.Println()
	t.PrintTableHeader("Stat", "Value")
	t.PrintTableLine("Avg Latency", fmt.Sprintf("%.2fs", avgLatency))
	t.PrintTableLine("Max Latency", fmt.Sprintf("%.2fs", t.MaxLatency.Seconds()))
	t.PrintTableLine("Req / Sec", fmt.Sprintf("%.2f", messagesPerSecond))
	t.PrintTableFooter()
	t.Println()
}

func (t *Tester) printErrors() {
	var errorCount int
	for key, messages := range t.messages {
		if strings.HasSuffix(key, "Error") {
			for _, message := range messages {
				if errorCount == 0 {
					t.Println(HeaderColor.Sprint("\nErrors"))
					t.Println("===========================")
				}

				t.Println(key, ": ", ErrorColor.Sprint(message))
				errorCount++
			}
		}
	}

}

func (t *Tester) PrintInt(key string, value int) {
	t.PrintValue(key, fmt.Sprint(value))
}

func (t *Tester) PrintDuration(key string, value time.Duration) {
	t.PrintValue(key, value.Round(50*time.Millisecond).String())
}

func (t *Tester) PrintValue(key, value string) {
	t.Println(StatColor.Sprintf("%s:\t%s", key, value))
}

func (t *Tester) PrintTableHeader(header1, header2 string) {
	t.Println(LineColor.Sprint("┌───────────────┬────────────┐"))
	t.Println(t.pipe(), HeaderColor.Sprintf(" %s\t\t", header1), t.pipe(), HeaderColor.Sprintf(" %s\t     ", header2), t.pipe())
	t.Println(LineColor.Sprint("├───────────────┼────────────┤"))
}

func (t *Tester) PrintTableLine(value1, value2 string) {
	t.Println(
		t.pipe(),
		StatColor.Sprintf(" %s%s", value1, spacesBuffer(value1, 14)),
		t.pipe(),
		StatColor.Sprintf(" %s%s", value2, spacesBuffer(value2, 11)),
		t.pipe(),
	)
}

// spacesBuffer returns the number of spaces needed to fill the remaining space in the table.
func spacesBuffer(value string, size int) (spaces string) {
	for i := 0; i < size-len(value); i++ {
		spaces += " "
	}
	return spaces
}

func (t *Tester) pipe() string {
	return LineColor.Sprint("│")
}

func (t *Tester) PrintTableFooter() {
	t.Println(LineColor.Sprint("└───────────────┴────────────┘"))
}

func (t *Tester) Printf(format string, msgs ...interface{}) {
	t.Config.PlatformFixture.Printf(format, msgs...)
}

func (t *Tester) Println(msgs ...interface{}) {
	t.Config.PlatformFixture.Println(msgs...)
}

func (t *Tester) Print(msgs ...interface{}) {
	t.Config.PlatformFixture.Print(msgs)
}
