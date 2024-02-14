package allocator

import (
	"encoding/json"
	"fmt"
	platform "github.com/PixoVR/pixo-golang-clients/pixo-platform/primary-api"
	"github.com/rs/zerolog/log"
)

func (a *AllocatorClient) RegisterTrigger(trigger platform.MultiplayerServerTrigger) Response {

	body, err := json.Marshal(trigger)
	if err != nil {
		log.Debug().Err(err).Msg("Failed to marshal trigger registration request")
		return Response{Error: err}
	}

	path := "allocator/build/triggers"

	res, err := a.Post(path, body)
	if err != nil {
		log.Debug().Err(err).Msg("Failed to post trigger registration request")
		return Response{
			HTTPResponse: res,
			Error:        err,
		}
	}

	return Response{HTTPResponse: res}
}

func (a *AllocatorClient) UpdateTrigger(trigger platform.MultiplayerServerTrigger) Response {

	body, err := json.Marshal(trigger)
	if err != nil {
		log.Debug().Err(err).Msg("Failed to marshal trigger update request")
		return Response{Error: err}
	}

	path := fmt.Sprintf("allocator/build/triggers/%d", trigger.ID)

	res, err := a.Put(path, body)
	if err != nil {
		log.Debug().Err(err).Msg("Failed to put trigger update request")
		return Response{
			HTTPResponse: res,
			Error:        err,
		}
	}

	return Response{HTTPResponse: res}
}

func (a *AllocatorClient) DeleteTrigger(id int) Response {

	path := fmt.Sprintf("allocator/build/triggers/%d", id)

	res, err := a.Delete(path)
	if err != nil {
		log.Debug().Err(err).Msg("Failed to delete trigger")
		return Response{
			HTTPResponse: res,
			Error:        err,
		}
	}

	return Response{HTTPResponse: res}
}
