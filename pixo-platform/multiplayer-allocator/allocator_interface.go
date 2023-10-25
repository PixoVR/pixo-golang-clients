package multiplayer_allocator

import (
	"github.com/go-resty/resty/v2"
)

type Allocator interface {
	AllocateGameserver(request AllocationRequest) AllocationResponse
	RegisterFleet(fleet FleetRequest) (*resty.Response, error)
}
