package matchmaker

import (
	"net"
	"time"
)

// DialGameserver connects to the gameserver address given
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

// CloseGameserverConnection closes the connection to the gameserver
func (m *MultiplayerMatchmaker) CloseGameserverConnection() error {
	if err := m.gameserverConnection.Close(); err != nil {
		return err
	}

	return nil
}

// SendMessageToGameserver sends a message to the gameserver
func (m *MultiplayerMatchmaker) SendMessageToGameserver(message []byte) error {
	if _, err := m.gameserverConnection.Write(message); err != nil {
		return err
	}

	return nil
}

// ReadMessageFromGameserver reads a message from the gameserver
func (m *MultiplayerMatchmaker) ReadMessageFromGameserver() ([]byte, error) {
	if m.gameserverConnection == nil {
		if err := m.DialGameserver(m.gameserverAddress); err != nil {
			return nil, err
		}
	}

	if err := m.gameserverConnection.SetReadDeadline(time.Now().Add(5 * time.Second)); err != nil {
		return nil, err
	}

	received := make([]byte, 1024)
	n, err := m.gameserverConnection.Read(received)
	if err != nil {
		return nil, err
	}

	response := received[:n]

	return response, nil
}
