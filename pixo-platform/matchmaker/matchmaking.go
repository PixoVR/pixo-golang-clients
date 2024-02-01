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
		return nil, errors.New("match request is invalid")
	}

	message, err := json.Marshal(req)
	if err != nil {
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
		return nil, err
	}

	if !matchResponse.IsValid() {
		return nil, errors.New(matchResponse.Message)
	}

	port, err := strconv.Atoi(matchResponse.MatchDetails.Port)
	if err != nil {
		return nil, err
	}

	m.gameserverAddress = &net.UDPAddr{
		IP:   net.ParseIP(matchResponse.MatchDetails.IP),
		Port: port,
	}
	return m.gameserverAddress, nil
}

func (m *MultiplayerMatchmaker) DialGameserver(addr *net.UDPAddr) error {

	udpServer, err := net.ResolveUDPAddr(addr.Network(), addr.String())
	if err != nil {
		return err
	}

	conn, err := net.DialUDP(addr.Network(), nil, udpServer)
	if err != nil {
		return err
	}

	m.gameserverConnection = conn
	return nil
}

func (m *MultiplayerMatchmaker) SendAndReceiveMessage(message []byte) ([]byte, error) {
	if m.gameserverConnection == nil {
		return nil, errors.New("gameserver connection is nil")
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
		return err
	}

	return nil
}

func (m *MultiplayerMatchmaker) sendGameServerMessage(message []byte) error {
	if _, err := m.gameserverConnection.Write(message); err != nil {
		return err
	}

	return nil
}

func (m *MultiplayerMatchmaker) readGameServerMessage() ([]byte, error) {
	received := make([]byte, 1024)
	n, err := m.gameserverConnection.Read(received)
	if err != nil {
		return nil, err
	}

	response := received[:n]

	return response, nil
}
