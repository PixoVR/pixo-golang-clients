package matchmaker

import (
	"encoding/json"
	"errors"
	"io"
	"net"
	"strconv"
)

func (m *MultiplayerMatchmaker) FindMatch(req MatchRequest) (*net.UDPAddr, error) {

	if !req.IsValid() {
		return nil, errors.New("match request is invalid")
	}

	httpResponse, err := m.DialWebsocket(MatchmakingEndpoint)
	if err != nil {
		return nil, err
	}

	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(httpResponse.Body)

	message, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	if err = m.WriteToWebsocket(message); err != nil {
		return nil, err
	}

	_, response, err := m.ReadFromWebsocket()
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
