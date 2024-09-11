package allocator

type Allocator interface {
	AllocateGameserver(request AllocationRequest) AllocationResponse
}
