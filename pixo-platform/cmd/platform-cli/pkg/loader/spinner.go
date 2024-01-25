package loader

import (
	"fmt"
	"io"
	"time"
)

type Spinner struct {
	periodSeconds time.Duration
	chars         []rune
	writer        io.Writer
	doneChan      chan bool
}

func NewSpinner(writer io.Writer) *Spinner {

	period := 250 * time.Millisecond

	spinner := &Spinner{
		periodSeconds: period,
		writer:        writer,
		chars:         []rune{'-', '\\', '|', '/'},
		doneChan:      make(chan bool),
	}

	go spinner.Start()

	return spinner
}

func (s *Spinner) Start() {
	for {
		for _, r := range s.chars {
			select {
			case <-s.doneChan:
				return
			default:
				s.writer.Write([]byte(fmt.Sprintf("\r%c", r))) // nolint: errcheck
				time.Sleep(s.period())
			}
		}
	}
}

func (s *Spinner) Stop() {
	s.writer.Write([]byte("\r")) // nolint: errcheck
	s.doneChan <- true
	close(s.doneChan)
}

func (s *Spinner) period() time.Duration {
	return s.periodSeconds / time.Duration(len(s.chars))
}
