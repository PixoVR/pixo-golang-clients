package allocator

import (
	"fmt"
)

// Spy is a struct that contains fields to track the number of times a method is called.
type Spy struct {
	CalledAllocateGameserver bool
}

// NewAllocatorSpy creates a new instance of the Spy struct.
func NewAllocatorSpy() *Spy {
	return &Spy{}
}

// AllocateGameserver is a spy method that sets the CalledAllocateGameserver field to true and returns a dummy AllocationResponse.
func (a *Spy) AllocateGameserver(request AllocationRequest) AllocationResponse {
	a.CalledAllocateGameserver = true

	return AllocationResponse{
		Results: GameServer{
			Name: "test-gameserver",
			IP:   Localhost,
			Port: fmt.Sprint(DefaultGameserverPort),
		},
	}
}
