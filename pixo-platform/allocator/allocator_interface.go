package allocator

import (
	"github.com/PixoVR/pixo-golang-clients/pixo-platform/platform"
)

type Allocator interface {
	AllocateGameserver(request AllocationRequest) AllocationResponse
	RegisterFleet(fleet FleetRequest) Response
	DeregisterFleet(fleet FleetRequest) Response
	RegisterTrigger(trigger platform.MultiplayerServerTrigger) Response
	UpdateTrigger(trigger platform.MultiplayerServerTrigger) Response
	DeleteTrigger(id int) Response
}
