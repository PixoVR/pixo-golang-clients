package matchmaker

import (
	"encoding/json"
	"errors"
	"fmt"
	abstractClient "github.com/PixoVR/pixo-golang-clients/pixo-platform/abstract-client"
	"github.com/rs/zerolog/log"
	"io"
	"net"
	"strconv"
)

type MultiplayerMatchmaker struct {
	abstractClient.PixoAbstractAPIClient
}

func NewMatchmaker(url, token string) *MultiplayerMatchmaker {

	if url == "" {
		url = getURL()
	}

	return &MultiplayerMatchmaker{
		PixoAbstractAPIClient: *abstractClient.NewClient(token, url),
	}
}

func getURL() string {
	return fmt.Sprintf("%s/%s", DefaultMatchmakingURL, MatchmakingEndpoint)
}

func (p *MultiplayerMatchmaker) Connect(moduleID, orgID int) (*net.UDPAddr, error) {
	log.Info().Msg("Connecting to matchmaking server")

	httpResponse, err := p.ConnectToWebsocket()
	if err != nil {
		return nil, err
	}

	defer func(Body io.ReadCloser) {
		if err = Body.Close(); err != nil {
			log.Error().Err(err).Msg("Error closing HTTP response body")
		}
	}(httpResponse.Body)

	match := MatchRequest{
		ModuleID: moduleID,
		OrgID:    orgID,
	}

	if !match.IsValid() {
		err = errors.New("match request is invalid")
		log.Error().Err(err).Msg("Match request is invalid")
		return nil, err
	}

	message, err := json.Marshal(match)
	if err != nil {
		log.Error().Err(err).Msg("Error deserializing match request")
		return nil, err
	}

	if err = p.SendMessageToWebsocket(message); err != nil {
		return nil, err
	}

	response, err := p.ReadFromWebsocket()
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

	return &net.UDPAddr{
		IP:   net.ParseIP(matchResponse.MatchDetails.IP),
		Port: port,
	}, nil
}
