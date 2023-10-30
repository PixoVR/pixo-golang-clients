package multiplayer_allocator

import (
	"fmt"
	"github.com/PixoVR/pixo-golang-clients/pixo-platform/matchmaker"
	platform "github.com/PixoVR/pixo-golang-clients/pixo-platform/primary-api"
	"github.com/go-resty/resty/v2"
	"net/http"
)

type AllocatorSpy struct {
	CalledAllocateGameserver bool
	CalledRegisterFleet      bool
	CalledDeregisterFleet    bool
	CalledRegisterTrigger    bool
	CalledUpdateTrigger      bool
	CalledDeleteTrigger      bool
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
		Results: GameServer{
			ResourceName: "test-gameserver",
			Address:      matchmaker.Localhost,
			Port:         fmt.Sprint(matchmaker.DefaultGameserverPort),
		},
	}
}

func (a *AllocatorSpy) RegisterFleet(fleet FleetRequest) Response {
	a.CalledRegisterFleet = true

	return Response{
		HTTPResponse: &resty.Response{
			RawResponse: &http.Response{
				StatusCode: http.StatusCreated,
			},
		},
	}
}

func (a *AllocatorSpy) DeregisterFleet(fleet FleetRequest) Response {
	a.CalledDeregisterFleet = true

	return Response{
		HTTPResponse: &resty.Response{
			RawResponse: &http.Response{
				StatusCode: http.StatusNoContent,
			},
		},
	}
}

func (b *AllocatorSpy) RegisterTrigger(trigger platform.MultiplayerServerTrigger) Response {
	b.CalledRegisterTrigger = true

	return Response{
		HTTPResponse: &resty.Response{
			RawResponse: &http.Response{
				StatusCode: http.StatusCreated,
			},
		},
	}
}

func (b *AllocatorSpy) UpdateTrigger(trigger platform.MultiplayerServerTrigger) Response {
	b.CalledUpdateTrigger = true

	return Response{
		HTTPResponse: &resty.Response{
			RawResponse: &http.Response{
				StatusCode: http.StatusOK,
			},
		},
	}
}

func (b *AllocatorSpy) DeleteTrigger(id int) Response {
	b.CalledDeleteTrigger = true

	return Response{
		HTTPResponse: &resty.Response{
			RawResponse: &http.Response{
				StatusCode: http.StatusNoContent,
			},
		},
	}
}
