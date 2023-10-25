package multiplayer_allocator

import (
	"github.com/go-resty/resty/v2"
	"net/http"
)

type AllocatorSpy struct {
	CalledAllocateGameserver bool
	CalledRegisterFleet      bool
	CalledDeregisterFleet    bool
}

func NewAllocatorSpy() *AllocatorSpy {
	return &AllocatorSpy{
		CalledAllocateGameserver: false,
		CalledRegisterFleet:      false,
		CalledDeregisterFleet:    false,
	}
}

func (a *AllocatorSpy) AllocateGameserver(request AllocationRequest) AllocationResponse {
	a.CalledAllocateGameserver = true

	return AllocationResponse{
		HTTPResponse: &resty.Response{
			RawResponse: &http.Response{
				StatusCode: http.StatusCreated,
			},
		},
	}
}

func (a *AllocatorSpy) RegisterFleet(fleet FleetRequest) (*resty.Response, error) {
	a.CalledRegisterFleet = true

	return &resty.Response{
		RawResponse: &http.Response{
			StatusCode: http.StatusCreated,
		},
	}, nil
}

func (a *AllocatorSpy) DeregisterFleet(fleet FleetRequest) (*resty.Response, error) {
	a.CalledDeregisterFleet = true

	return &resty.Response{
		RawResponse: &http.Response{
			StatusCode: http.StatusNoContent,
		},
	}, nil
}
