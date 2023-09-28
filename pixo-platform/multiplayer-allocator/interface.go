package multiplayer_allocator

import (
	platform "github.com/PixoVR/pixo-golang-clients/pixo-platform/primary-api"
	"github.com/go-resty/resty/v2"
	"net/http"
)

type Allocator interface {
	AllocateGameserver(request AllocationRequest) (AllocationResponse, error)
	RegisterTrigger(trigger platform.MultiplayerServerTrigger) (*resty.Response, error)
	DeleteTrigger(id int) (*resty.Response, error)
}

type AllocatorSpy struct {
	CalledAllocateGameserver bool
	CalledRegisterTrigger    bool
	CalledDeleteTrigger      bool
}

func NewAllocatorSpy() *AllocatorSpy {
	return &AllocatorSpy{}
}

func (a *AllocatorSpy) AllocateGameserver(request AllocationRequest) (AllocationResponse, error) {
	a.CalledAllocateGameserver = true

	return AllocationResponse{
		HTTPResponse: &resty.Response{
			RawResponse: &http.Response{
				StatusCode: http.StatusCreated,
			},
		},
	}, nil
}

func (a *AllocatorSpy) RegisterTrigger(trigger platform.MultiplayerServerTrigger) (*resty.Response, error) {
	a.CalledRegisterTrigger = true

	return &resty.Response{
		RawResponse: &http.Response{
			StatusCode: http.StatusCreated,
		},
	}, nil
}

func (a *AllocatorSpy) DeleteTrigger(id int) (*resty.Response, error) {
	a.CalledDeleteTrigger = true

	return &resty.Response{
		RawResponse: &http.Response{
			StatusCode: http.StatusNoContent,
		},
	}, nil
}
