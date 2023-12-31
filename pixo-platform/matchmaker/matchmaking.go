package matchmaker

import (
	"encoding/json"
	"errors"
	"github.com/rs/zerolog/log"
	"io"
	"net"
	"strconv"
)

func (m *MultiplayerMatchmaker) FindMatch(req MatchRequest) (*net.UDPAddr, error) {
	log.Debug().Msg("Connecting to matchmaking server")

	httpResponse, err := m.ConnectToWebsocket()
	if err != nil {
		return nil, err
	}

	defer func(Body io.ReadCloser) {
		if err = Body.Close(); err != nil {
			log.Error().Err(err).Msg("Error closing HTTP response body")
		}
	}(httpResponse.Body)

	if !req.IsValid() {
		err = errors.New("match request is invalid")
		log.Error().Err(err).Msg("Match request is invalid")
		return nil, err
	}

	message, err := json.Marshal(req)
	if err != nil {
		log.Error().Err(err).Msg("Error deserializing match request")
		return nil, err
	}

	if err = m.SendMessageToWebsocket(message); err != nil {
		return nil, err
	}

	response, err := m.ReadFromWebsocket()
	if err != nil {
		return nil, err
	}

	var matchResponse MatchResponse
	if err = json.Unmarshal(response, &matchResponse); err != nil {
		log.Error().Err(err).Msg("Error serializing match response")
		return nil, err
	}

	if !matchResponse.IsValid() {
		err = errors.New("match response is invalid")
		log.Error().Err(err)
		return nil, err
	}

	port, err := strconv.Atoi(matchResponse.MatchDetails.Port)
	if err != nil {
		log.Error().Err(err).Msg("Error parsing port from match response")
		return nil, err
	}

	m.gameserverAddress = &net.UDPAddr{
		IP:   net.ParseIP(matchResponse.MatchDetails.IP),
		Port: port,
	}
	return m.gameserverAddress, nil
}

func (m *MultiplayerMatchmaker) DialGameserver(addr *net.UDPAddr) error {
	log.Debug().Msg("Connecting to gameserver")

	udpServer, err := net.ResolveUDPAddr(addr.Network(), addr.String())
	if err != nil {
		log.Error().Err(err).Msg("unable to resolve address")
		return err
	}

	conn, err := net.DialUDP(addr.Network(), nil, udpServer)
	if err != nil {
		log.Error().Err(err).Msg("unable to dial gameserver address")
		return err
	}

	m.gameserverConnection = conn
	return nil
}

func (m *MultiplayerMatchmaker) SendAndReceiveMessage(message []byte) ([]byte, error) {
	if m.gameserverConnection == nil {
		err := errors.New("gameserver connection is nil")
		log.Debug().Err(err).Msg("unable to send message to gameserver")
		return nil, err
	}

	if err := m.sendGameServerMessage(message); err != nil {
		return nil, err
	}

	response, err := m.readGameServerMessage()
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (m *MultiplayerMatchmaker) CloseGameserverConnection() error {
	if err := m.gameserverConnection.Close(); err != nil {
		log.Error().Err(err).Msg("unable to close gameserver connection")
		return err
	}

	log.Debug().Msg("Closed gameserver connection")
	return nil
}

func (m *MultiplayerMatchmaker) sendGameServerMessage(message []byte) error {
	if _, err := m.gameserverConnection.Write(message); err != nil {
		log.Error().Err(err).Msg("unable to write to gameserver")
		return err
	}

	log.Debug().Msgf("Sent message to gameserver: %s", message)
	return nil
}

func (m *MultiplayerMatchmaker) readGameServerMessage() ([]byte, error) {
	received := make([]byte, 1024)
	n, err := m.gameserverConnection.Read(received)
	if err != nil {
		log.Error().Err(err).Msg("unable to read from gameserver")
		return nil, err
	}

	response := received[:n]

	log.Debug().Msgf("Received message from gameserver: %s", response)
	return response, nil
}
