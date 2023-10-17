package multiplayer_allocator

import (
	"encoding/json"
	"github.com/rs/zerolog/log"
)

func (a *AllocatorClient) RegisterFleet(request FleetRegisterRequest) FleetRegisterResponse {

	body, err := json.Marshal(request)
	if err != nil {
		log.Debug().Err(err).Msg("Failed to marshal register fleet request")
		return FleetRegisterResponse{
			Error: err,
		}
	}

	path := "allocator/register"

	res, err := a.Post(path, body)
	if err != nil {
		log.Debug().Err(err).Msg("Failed to post register fleet request")
		return FleetRegisterResponse{
			HTTPResponse: res,
			Error:        err,
		}
	}

	var response FleetRegisterResponse
	err = json.Unmarshal(res.Body(), &response)
	if err != nil {
		log.Debug().Err(err).Msg("Failed to unmarshal register fleet response")
		response.Error = err
	}

	response.HTTPResponse = res

	return response
}
