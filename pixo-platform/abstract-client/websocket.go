package abstract_client

import (
	"errors"
	"github.com/gorilla/websocket"
	"github.com/rs/zerolog/log"
	"net/http"
	"time"
)

func (p *PixoAbstractAPIClient) ConnectToWebsocket() (*http.Response, error) {
	log.Info().Msg("Connecting to websocket")

	httpHeader := http.Header{}
	if p.token != "" {
		httpHeader.Add("Authorization", "Bearer "+p.token)
	}

	conn, httpResponse, err := websocket.DefaultDialer.Dial(p.url, httpHeader)
	if err != nil {
		log.Error().Err(err).Msg("Error connecting to websocket")
		return httpResponse, err
	}

	if err = conn.SetReadDeadline(time.Now().Add(p.timeoutDuration())); err != nil {
		log.Info().Err(err).Msg("Error setting read deadline for websocket")
		return httpResponse, err
	}

	p.conn = conn

	log.Info().Msg("Successfully connected to websocket")
	return httpResponse, nil
}

func (p *PixoAbstractAPIClient) SendMessageToWebsocket(message []byte) error {
	log.Info().Msg("Sending message to websocket")

	if err := p.conn.WriteMessage(websocket.TextMessage, message); err != nil {
		log.Error().Err(err).Msg("Error during writing to websocket")
		return err
	}

	log.Info().Msgf("Sent message to websocket: %s", message)
	return nil
}

func (p *PixoAbstractAPIClient) ReadFromWebsocket() ([]byte, error) {
	log.Info().Msg("Reading message from websocket")

	if err := p.conn.SetReadDeadline(time.Now().Add(p.timeoutDuration())); err != nil {
		log.Info().Err(err).Msg("Error setting read deadline")
		return nil, err
	}

	_, msg, err := p.conn.ReadMessage()
	if err != nil {
		log.Error().Err(err).Msg("Error reading message from websocket")
		return nil, err
	}

	if len(msg) == 0 {
		return nil, errors.New("empty response data")
	}

	log.Info().Msgf("Received message from websocket: %s", msg)
	return msg, nil
}

func (p *PixoAbstractAPIClient) CloseWebsocketConnection() error {
	log.Info().Msg("Closing websocket")

	if err := p.conn.Close(); err != nil {
		log.Error().Err(err).Msg("Error closing websocket")
		return err
	}

	log.Info().Msg("Websocket closed")
	return nil
}

func (p *PixoAbstractAPIClient) timeoutDuration() time.Duration {
	if p.timeoutSeconds == 0 {
		return 30 * time.Second
	}

	return time.Duration(p.timeoutSeconds) * time.Second
}
