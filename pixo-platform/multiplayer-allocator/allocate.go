package multiplayer_allocator

import (
	"encoding/json"
	"github.com/rs/zerolog/log"
)

func (p *AllocatorClient) AllocateGameserver(request AllocationRequest) (AllocationResponse, error) {

	body, err := json.Marshal(request)
	if err != nil {
		log.Debug().Err(err).Msg("Failed to marshal allocate server request")
		return AllocationResponse{}, err
	}

	path := "allocator/allocate"

	res, err := p.Post(path, body)
	if err != nil {
		log.Debug().Err(err).Msg("Failed to post allocate server request")
		return AllocationResponse{}, err
	}

	var response AllocationResponse
	err = json.Unmarshal(res.Body(), &response)
	if err != nil {
		log.Debug().Err(err).Msg("Failed to unmarshal allocate server response")
		return response, err
	}

	response.HTTPResponse = res

	return response, nil
}
