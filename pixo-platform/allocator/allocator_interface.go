package allocator

import (
	platform "github.com/PixoVR/pixo-golang-clients/pixo-platform/legacy"
)

type Allocator interface {
	AllocateGameserver(request AllocationRequest) AllocationResponse
	RegisterFleet(fleet FleetRequest) Response
	DeregisterFleet(fleet FleetRequest) Response
	RegisterTrigger(trigger platform.MultiplayerServerTrigger) Response
	UpdateTrigger(trigger platform.MultiplayerServerTrigger) Response
	DeleteTrigger(id int) Response
}
