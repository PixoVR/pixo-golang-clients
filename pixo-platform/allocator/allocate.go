package allocator

import (
	"encoding/json"
	"github.com/rs/zerolog/log"
)

func (a *Client) AllocateGameserver(request AllocationRequest) AllocationResponse {

	body, err := json.Marshal(request)
	if err != nil {
		log.Debug().Err(err).Msg("Failed to marshal allocate server request")
		return AllocationResponse{Error: err}
	}

	path := "allocate"

	res, err := a.Post(path, body)
	if err != nil {
		log.Debug().Err(err).Msg("Failed to post allocate server request")
		return AllocationResponse{
			Error: err,
		}
	}

	var response AllocationResponse
	if err = json.Unmarshal(res.Body(), &response); err != nil {
		log.Debug().Err(err).Msg("Failed to unmarshal allocate server response")
		response.Error = err
		return response
	}

	return response
}
