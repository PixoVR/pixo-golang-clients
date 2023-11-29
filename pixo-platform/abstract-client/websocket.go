package abstract_client

import (
	"errors"
	"github.com/gorilla/websocket"
	"github.com/rs/zerolog/log"
	"net/http"
	"time"
)

func (p *PixoAbstractAPIClient) ConnectToWebsocket() (*http.Response, error) {
	log.Debug().Msg("Connecting to websocket")

	httpHeader := http.Header{}
	if p.token != "" {
		httpHeader.Add("Authorization", "Bearer "+p.token)
	}

	conn, httpResponse, err := websocket.DefaultDialer.Dial(p.url, httpHeader)
	if err != nil {
		log.Error().Err(err).Msg("unable to connect to websocket")
		return nil, err
	}

	if err = conn.SetReadDeadline(time.Now().Add(p.timeoutDuration())); err != nil {
		log.Error().Err(err).Msg("unable to set read deadline for websocket")
		return nil, err
	}

	p.conn = conn

	log.Debug().Msg("Successfully connected to websocket")
	return httpResponse, nil
}

func (p *PixoAbstractAPIClient) SendMessageToWebsocket(message []byte) error {
	log.Debug().Msg("Sending message to websocket")

	if err := p.conn.WriteMessage(websocket.TextMessage, message); err != nil {
		log.Error().Err(err).Msg("unable to write to websocket")
		return err
	}

	log.Debug().Msgf("Sent message to websocket: %s", message)
	return nil
}

func (p *PixoAbstractAPIClient) ReadFromWebsocket() ([]byte, error) {
	log.Debug().Msg("Reading message from websocket")

	_, msg, err := p.conn.ReadMessage()
	if err != nil {
		log.Error().Err(err).Msg("Error reading message from websocket")
		return nil, err
	}

	if len(msg) == 0 {
		return nil, errors.New("empty response data")
	}

	log.Debug().Msgf("Received message from websocket: %s", msg)
	return msg, nil
}

func (p *PixoAbstractAPIClient) CloseWebsocketConnection() error {
	log.Debug().Msg("Closing websocket")

	if err := p.conn.Close(); err != nil {
		log.Error().Err(err).Msg("Error closing websocket")
		return err
	}

	log.Debug().Msg("Websocket closed")
	return nil
}

func (p *PixoAbstractAPIClient) timeoutDuration() time.Duration {
	if p.timeoutSeconds == 0 {
		return 30 * time.Second
	}

	return time.Duration(p.timeoutSeconds) * time.Second
}
