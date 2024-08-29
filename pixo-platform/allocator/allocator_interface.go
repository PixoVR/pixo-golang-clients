package allocator

import (
	"github.com/PixoVR/pixo-golang-clients/pixo-platform/legacy"
)

type Allocator interface {
	AllocateGameserver(request AllocationRequest) AllocationResponse
	RegisterFleet(fleet FleetRequest) Response
	DeregisterFleet(fleet FleetRequest) Response
	RegisterTrigger(trigger legacy.MultiplayerServerTrigger) Response
	UpdateTrigger(trigger legacy.MultiplayerServerTrigger) Response
	DeleteTrigger(id int) Response
}
