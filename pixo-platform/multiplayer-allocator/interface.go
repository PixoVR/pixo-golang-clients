package multiplayer_allocator

import (
	platform "github.com/PixoVR/pixo-golang-clients/pixo-platform/primary-api"
	"github.com/go-resty/resty/v2"
)

type Allocator interface {
	AllocateGameserver(request AllocationRequest) AllocationResponse
	RegisterTrigger(trigger platform.MultiplayerServerTrigger) (*resty.Response, error)
	DeleteTrigger(id int) (*resty.Response, error)
}
