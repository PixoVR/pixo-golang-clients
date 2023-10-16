package multiplayer_allocator

import (
	platform "github.com/PixoVR/pixo-golang-clients/pixo-platform/primary-api"
	"github.com/go-resty/resty/v2"
	"net/http"
)

type AllocatorSpy struct {
	CalledAllocateGameserver bool
	CalledRegisterTrigger    bool
	CalledUpdateTrigger      bool
	CalledDeleteTrigger      bool
}

func NewAllocatorSpy() *AllocatorSpy {
	return &AllocatorSpy{
		CalledAllocateGameserver: false,
		CalledRegisterTrigger:    false,
		CalledUpdateTrigger:      false,
		CalledDeleteTrigger:      false,
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

func (a *AllocatorSpy) RegisterTrigger(trigger platform.MultiplayerServerTrigger) (*resty.Response, error) {
	a.CalledRegisterTrigger = true

	return &resty.Response{
		RawResponse: &http.Response{
			StatusCode: http.StatusCreated,
		},
	}, nil
}

func (a *AllocatorSpy) UpdateTrigger(trigger platform.MultiplayerServerTrigger) (*resty.Response, error) {
	a.CalledUpdateTrigger = true

	return &resty.Response{
		RawResponse: &http.Response{
			StatusCode: http.StatusOK,
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
