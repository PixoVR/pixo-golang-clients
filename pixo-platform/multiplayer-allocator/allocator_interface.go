package multiplayer_allocator

import (
	"github.com/go-resty/resty/v2"
)

type ServerAllocatorClient interface {
	AllocateGameserver(request AllocationRequest) AllocationResponse
	RegisterFleet(fleet FleetRequest) (*resty.Response, error)
}
