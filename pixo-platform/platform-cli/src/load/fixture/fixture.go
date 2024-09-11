package fixture

import (
	"github.com/PixoVR/pixo-golang-clients/pixo-platform/platform-cli/src/ctx"
	"github.com/PixoVR/pixo-golang-clients/pixo-platform/platform-cli/src/printer"
	"github.com/spf13/cobra"
	"io"
	"os"
	"sync"
	"time"
)

// Config contains the configuration for a load test.
type Config struct {
	Command         *cobra.Command
	PlatformFixture *ctx.Context
	Writer          io.Writer
	Amount          int
	Concurrent      int
	MaxDuration     time.Duration
}

// Tester configures and runs WebSocket load tests.
type Tester struct {
	printer.Printer
	Config       Config
	NumLatencies int
	NumDone      int
	TotalLatency time.Duration
	MaxLatency   time.Duration
	start        time.Time
	end          time.Time
	messages     map[string][]string
	taskQueue    chan int
	wg           sync.WaitGroup
	lock         sync.Mutex
}

// NewLoadTester creates a new instance of Tester with the specified configuration.
func NewLoadTester(config Config) *Tester {
	t := &Tester{
		Config:   config,
		wg:       sync.WaitGroup{},
		messages: make(map[string][]string),
	}

	if t.Config.Writer == nil {
		t.Config.Writer = os.Stdout
	}
	t.Printer = printer.NewEmojiPrinter(t.Config.Writer)

	if t.Config.Amount <= 0 {
		amount, ok := t.Config.PlatformFixture.ConfigManager.GetIntFlagOrConfigValue("amount", t.Config.Command)
		if ok {
			t.Config.Amount = amount
		} else {
			t.Config.Amount = 50
		}
	}

	if t.Config.Concurrent <= 0 {
		concurrent, ok := t.Config.PlatformFixture.ConfigManager.GetIntFlagOrConfigValue("concurrent", t.Config.Command)
		if ok {
			t.Config.Concurrent = concurrent
		} else {
			t.Config.Concurrent = 10
		}
	}
	t.Config.Concurrent = min(t.Config.Concurrent, t.Config.Amount)
	t.taskQueue = make(chan int, t.Config.Concurrent)

	if t.Config.MaxDuration <= 0 {
		timeout, ok := t.Config.PlatformFixture.ConfigManager.GetIntFlagOrConfigValue("timeout", t.Config.Command)
		if ok {
			t.Config.MaxDuration = time.Duration(timeout) * time.Minute
		} else {
			t.Config.MaxDuration = 2 * time.Minute
		}
	}

	return t
}

// Run starts the load testing process.
func (t *Tester) Run(performFunc func(id int)) {
	t.Println()
	t.Printf(":rocket: Starting load test with %d requests and %d concurrent workers\n\n", t.Config.Amount, t.Config.Concurrent)

	t.startWorkers(performFunc)

	t.start = time.Now()
	for i := 0; i < t.Config.Amount; i++ {
		t.taskQueue <- i + 1
	}
	close(t.taskQueue)
	t.wg.Wait()
	t.end = time.Now()

	t.DisplaySummary()
}

func (t *Tester) startWorkers(performFunc func(id int)) {
	for i := 0; i < t.Config.Concurrent; i++ {
		t.startWorker(performFunc)
	}
}

func (t *Tester) startWorker(performFunc func(id int)) {
	t.wg.Add(1)
	go func() {
		defer t.wg.Done()
		for i := range t.taskQueue {
			start := time.Now()
			performFunc(i)
			t.recordLatencySince(start)
		}
	}()
}
