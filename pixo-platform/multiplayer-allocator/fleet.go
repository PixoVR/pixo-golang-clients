package multiplayer_allocator

import (
	"encoding/json"
	"fmt"
	"github.com/rs/zerolog/log"
	"strings"
)

func (a *AllocatorClient) RegisterFleet(request FleetRequest) Response {

	body, err := json.Marshal(request)
	if err != nil {
		log.Debug().Err(err).Msg("Failed to marshal register fleet request")
		return Response{
			Error: err,
		}
	}

	path := "allocator/fleets"

	res, err := a.Post(path, body)
	if err != nil {
		log.Debug().Err(err).Msg("Failed to post register fleet request")
		return Response{
			HTTPResponse: res,
			Error:        err,
		}
	}

	var response Response
	err = json.Unmarshal(res.Body(), &response)
	if err != nil {
		log.Debug().Err(err).Msg("Failed to unmarshal register fleet response")
		response.Error = err
	}

	response.HTTPResponse = res

	return response
}

func (a *AllocatorClient) DeregisterFleet(request FleetRequest) Response {
	cleanedServerVersion := strings.ReplaceAll(request.ServerVersion.SemanticVersion, ".", "-")

	path := fmt.Sprintf("allocator/fleets/module/%d/serverVersion/%s", request.ServerVersion.ModuleID, cleanedServerVersion)

	res, err := a.Delete(path)
	if err != nil {
		log.Debug().Err(err).Msg("Failed to post deregister fleet request")
		return Response{
			HTTPResponse: res,
			Error:        err,
		}
	}

	var response Response
	if res.IsError() {
		log.Debug().Err(err).Msg("Failed to unmarshal deregister fleet response")
		response.Error = err
	}

	response.HTTPResponse = res

	return response
}
