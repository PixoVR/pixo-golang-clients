package allocator

import (
	"context"
	"encoding/json"
	"github.com/rs/zerolog/log"
	"io"
)

// AllocateGameserver send a request for a gameserver to the allocator
func (a *Client) AllocateGameserver(request AllocationRequest) (*GameServer, error) {
	body, err := json.Marshal(request)
	if err != nil {
		log.Error().Err(err).Msg("Failed to marshal allocate server request")
		return nil, err
	}

	path := "allocate"

	res, err := a.Post(context.TODO(), path, body)
	if err != nil {
		log.Error().Err(err).Msg("Failed to post allocate server request")
		return nil, err
	}

	resBody, _ := io.ReadAll(res.Body)

	var response AllocationResponse
	if err = json.Unmarshal(resBody, &response); err != nil {
		log.Error().Err(err).Msg("Failed to unmarshal allocate server response")
		return nil, err
	}

	return &response.Results, nil
}
