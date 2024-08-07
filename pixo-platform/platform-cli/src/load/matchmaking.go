package load

import (
	"errors"
	"github.com/fatih/color"
	"github.com/rs/zerolog/log"
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

// performRequest establishes a single WebSocket connection and requests a match
func (lt *Tester) performRequest(wg *sync.WaitGroup, id int) {
	defer wg.Done()

	log.Debug().Int("id", id).Msg("Calling dial")
	conn, _, err := lt.client.DialMatchmaker()
	if err != nil {
		lt.recordConnectionError(id, "failed to connect", err)
		return
	}
	defer func() {
		log.Debug().Int("id", id).Msg("Closing connection")
		_ = lt.client.CloseMatchmakerConnection(conn)
	}()

	lt.recordConnectionSuccess(id)

	start := time.Now()
	log.Debug().Int("id", id).Msg("Sending request")
	if err = lt.client.SendRequest(conn, lt.request); err != nil {
		lt.recordConnectionError(id, "failed to send request", err)
		return
	}

	lt.recordSentMessage()

	log.Debug().Int("id", id).Msg("Reading response")
	matchResponse, err := lt.client.ReadResponse(conn)
	lt.recordLatency(time.Since(start))
	lt.recordReceivedMessage()
	if err != nil {
		lt.recordConnectionError(id, "failed to read message", err)
		return
	}

	if !matchResponse.IsValid() {
		lt.recordMatchingError(id, "received invalid match", errors.New(matchResponse.Message))
		return
	}

	lt.recordSuccessMessageReceived(id, matchResponse)

	log.Debug().Int("id", id).Msg("Match received")
}
