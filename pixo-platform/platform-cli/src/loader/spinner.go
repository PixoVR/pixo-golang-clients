package loader

import (
	"context"
	"fmt"
	"io"
	"time"
)

type Loader struct {
	periodSeconds time.Duration
	chars         []rune
	writer        io.Writer
	doneChan      chan bool
}

func NewLoader(ctx context.Context, msg string, writer io.Writer) *Loader {

	period := 250 * time.Millisecond

	spinner := &Loader{
		periodSeconds: period,
		writer:        writer,
		chars:         []rune{'-', '\\', '|', '/'},
		doneChan:      make(chan bool),
	}

	//_, _ = writer.Write([]byte(emoji.Sprintf("%s ", msg)))

	go spinner.Start()

	return spinner
}

func (s *Loader) Start() {
	for {
		for _, r := range s.chars {
			select {
			case <-s.doneChan:
				return
			default:
				_, _ = s.writer.Write([]byte(fmt.Sprintf("\r%c", r)))
				time.Sleep(s.period())
			}
		}
	}
}

func (s *Loader) Stop() {
	_, _ = s.writer.Write([]byte("\b"))
	s.doneChan <- true
	close(s.doneChan)
}

func (s *Loader) period() time.Duration {
	return s.periodSeconds / time.Duration(len(s.chars))
}
