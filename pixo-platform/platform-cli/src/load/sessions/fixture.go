package sessions

import (
	"errors"
	"github.com/PixoVR/pixo-golang-clients/pixo-platform/headset"
	"github.com/PixoVR/pixo-golang-clients/pixo-platform/platform"
	"github.com/PixoVR/pixo-golang-clients/pixo-platform/platform-cli/src/load/fixture"
)

// Config contains the configuration for a load test.
type Config struct {
	fixture.Config
	Module       platform.Module
	Session      platform.Session
	EventRequest headset.EventRequest
	Legacy       bool
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

	config.Session = platform.Session{ModuleID: config.Module.ID, Module: config.Module}
	config.EventRequest = headset.EventRequest{ModuleID: config.Session.ModuleID}

	return &Tester{
		Tester: fixture.NewLoadTester(config.Config),
		config: config,
	}, nil
}

// Run starts the load testing process.
func (t *Tester) Run() {
	t.Tester.Run(t.performRequest)
	t.displayStats()
}
