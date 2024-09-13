package sessions

import (
	"encoding/json"
	"errors"
	"github.com/PixoVR/pixo-golang-clients/pixo-platform/platform"
	"github.com/PixoVR/pixo-golang-clients/pixo-platform/platform-cli/src/load/fixture"
)

// Config contains the configuration for a load test.
type Config struct {
	fixture.Config

	Session      platform.Session
	EventPayload string
	Event        platform.Event

	Legacy                 bool
	SessionStartDetails    map[string]interface{}
	SessionCompleteDetails map[string]interface{}
}

// Tester configures and runs sessions load tests.
type Tester struct {
	*fixture.Tester
	config Config
}

// NewLoadTester creates a new instance of Tester with the specified configuration.
func NewLoadTester(config Config) (*Tester, error) {
	if config.PlatformFixture == nil || config.PlatformFixture.PlatformClient == nil {
		return nil, errors.New("platform client is required")
	}

	var payload map[string]interface{}
	_ = json.Unmarshal([]byte(config.EventPayload), &payload)

	if config.Session.ModuleID == 0 {
		return nil, errors.New("module id is required")
	}

	config.Event = platform.Event{Payload: payload}

	t := &Tester{
		Tester: fixture.NewLoadTester(config.Config),
		config: config,
	}

	t.formatSessionStartDetails()
	t.formatSessionCompleteDetails()
	return t, nil
}

// Run starts the load testing process.
func (t *Tester) Run() {
	t.Tester.Run(t.performRequest)
	t.displayStats()
}
