package vertexai

import "github.com/rs/zerolog/log"

func (c *Client) debugMsg(msg string) {
	if c.debug {
		log.Debug().Msg(msg)
	}
}

func (c *Client) debugMsgf(format string, args ...interface{}) {
	if c.debug {
		log.Debug().Msgf(format, args...)
	}
}
