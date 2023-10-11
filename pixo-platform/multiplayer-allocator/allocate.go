package multiplayer_allocator

import (
	"encoding/json"
	"github.com/rs/zerolog/log"
)

func (a *AllocatorClient) AllocateGameserver(request AllocationRequest) AllocationResponse {

	body, err := json.Marshal(request)
	if err != nil {
		log.Debug().Err(err).Msg("Failed to marshal allocate server request")
		return AllocationResponse{Error: err}
	}

	path := "allocator/allocate"

	res, err := a.Post(path, body)
	if err != nil {
		log.Debug().Err(err).Msg("Failed to post allocate server request")
		return AllocationResponse{
			Error: err,
		}
	}

	var response AllocationResponse
	err = json.Unmarshal(res.Body(), &response)
	if err != nil {
		log.Debug().Err(err).Msg("Failed to unmarshal allocate server response")
		response.Error = err
		return response
	}

	response.HTTPResponse = res

	return response
}
