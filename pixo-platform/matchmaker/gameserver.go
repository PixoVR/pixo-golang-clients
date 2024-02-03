package matchmaker

import (
	"errors"
	"net"
)

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
