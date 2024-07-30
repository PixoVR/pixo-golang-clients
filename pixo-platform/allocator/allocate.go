package allocator

import (
	"encoding/json"
	"github.com/rs/zerolog/log"
)

func (a *Client) AllocateGameserver(request AllocationRequest) (*GameServer, error) {

	body, err := json.Marshal(request)
	if err != nil {
		log.Error().Err(err).Msg("Failed to marshal allocate server request")
		return nil, err
	}

	path := "allocate"

	res, err := a.Post(path, body)
	if err != nil {
		log.Error().Err(err).Msg("Failed to post allocate server request")
		return nil, err
	}

	var response AllocationResponse
	if err = json.Unmarshal(res.Body(), &response); err != nil {
		log.Error().Err(err).Msg("Failed to unmarshal allocate server response")
		return nil, err
	}

	return &response.Results, nil
}
