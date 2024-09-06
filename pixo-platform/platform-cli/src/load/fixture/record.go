package fixture

import (
	"fmt"
	"time"
)

// Messages returns the success messages received.
func (t *Tester) Messages(key string) []string {
	t.lock.Lock()
	defer t.lock.Unlock()

	return t.messages[key]
}

func (t *Tester) RecordSuccess(id int, key, msg string) {
	t.LogSuccess(id, msg)
	t.RecordMessage(key+"Success", msg)
}

func (t *Tester) RecordError(id int, key, msg string, err error) {
	t.LogError(id, msg, err)
	t.RecordMessage(key+"Error", fmt.Sprintf("%s: %s", msg, err.Error()))
}

// recordLatencySince tracks the latency of each request and updates max latency if necessary.
func (t *Tester) recordLatencySince(start time.Time) {
	t.lock.Lock()
	defer t.lock.Unlock()

	t.NumDone++
	latency := time.Since(start)

	t.TotalLatency += latency
	if latency > t.MaxLatency {
		t.MaxLatency = latency
	}

	t.NumLatencies++
}

// RecordMessage increments the count of successful messages received.
func (t *Tester) RecordMessage(key, msg string) {
	t.lock.Lock()
	defer t.lock.Unlock()

	t.messages[key] = append(t.messages[key], msg)
}
