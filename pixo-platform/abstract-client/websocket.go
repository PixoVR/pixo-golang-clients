package abstract_client

import (
	"errors"
	"github.com/gorilla/websocket"
	"github.com/rs/zerolog/log"
	"net/http"
	"time"
)

func (a *AbstractServiceClient) ConnectToWebsocket() (*http.Response, error) {
	log.Debug().Msg("Connecting to websocket")

	httpHeader := http.Header{}
	if a.token != "" {
		httpHeader.Add("Authorization", "Bearer "+a.token)
	}

	conn, httpResponse, err := websocket.DefaultDialer.Dial(a.url, httpHeader)
	if err != nil {
		log.Error().Err(err).Msg("unable to connect to websocket")
		return nil, err
	}

	if err = conn.SetReadDeadline(time.Now().Add(a.timeoutDuration())); err != nil {
		log.Error().Err(err).Msg("unable to set read deadline for websocket")
		return nil, err
	}

	a.conn = conn

	log.Debug().Msg("Successfully connected to websocket")
	return httpResponse, nil
}

func (a *AbstractServiceClient) SendMessageToWebsocket(message []byte) error {
	log.Debug().Msg("Sending message to websocket")

	if err := a.conn.WriteMessage(websocket.TextMessage, message); err != nil {
		log.Error().Err(err).Msg("unable to write to websocket")
		return err
	}

	log.Debug().Msgf("Sent message to websocket: %s", message)
	return nil
}

func (a *AbstractServiceClient) ReadFromWebsocket() ([]byte, error) {
	log.Debug().Msg("Reading message from websocket")

	_, msg, err := a.conn.ReadMessage()
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

func (a *AbstractServiceClient) CloseWebsocketConnection() error {
	log.Debug().Msg("Closing websocket")

	if err := a.conn.Close(); err != nil {
		log.Error().Err(err).Msg("Error closing websocket")
		return err
	}

	log.Debug().Msg("Websocket closed")
	return nil
}

func (a *AbstractServiceClient) timeoutDuration() time.Duration {
	if a.timeoutSeconds == 0 {
		return 30 * time.Second
	}

	return time.Duration(a.timeoutSeconds) * time.Second
}
