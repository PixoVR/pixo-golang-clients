package allocator

import (
	"fmt"
	platform "github.com/PixoVR/pixo-golang-clients/pixo-platform/legacy"
	"github.com/PixoVR/pixo-golang-clients/pixo-platform/matchmaker"
	"github.com/go-resty/resty/v2"
	"net/http"
)

type Spy struct {
	CalledAllocateGameserver bool
	CalledRegisterFleet      bool
	CalledDeregisterFleet    bool
	CalledRegisterTrigger    bool
	CalledUpdateTrigger      bool
	CalledDeleteTrigger      bool
}

func NewAllocatorSpy() *Spy {
	return &Spy{
		CalledAllocateGameserver: false,
		CalledRegisterFleet:      false,
		CalledDeregisterFleet:    false,
	}
}

func (a *Spy) AllocateGameserver(request AllocationRequest) AllocationResponse {
	a.CalledAllocateGameserver = true

	return AllocationResponse{
		Results: GameServer{
			Name: "test-gameserver",
			IP:   matchmaker.Localhost,
			Port: fmt.Sprint(matchmaker.DefaultGameserverPort),
		},
	}
}

func (a *Spy) RegisterFleet(fleet FleetRequest) Response {
	a.CalledRegisterFleet = true

	return Response{
		HTTPResponse: &resty.Response{
			RawResponse: &http.Response{
				StatusCode: http.StatusCreated,
			},
		},
	}
}

func (a *Spy) DeregisterFleet(fleet FleetRequest) Response {
	a.CalledDeregisterFleet = true

	return Response{
		HTTPResponse: &resty.Response{
			RawResponse: &http.Response{
				StatusCode: http.StatusNoContent,
			},
		},
	}
}

func (a *Spy) RegisterTrigger(trigger platform.MultiplayerServerTrigger) Response {
	a.CalledRegisterTrigger = true

	return Response{
		HTTPResponse: &resty.Response{
			RawResponse: &http.Response{
				StatusCode: http.StatusCreated,
			},
		},
	}
}

func (a *Spy) UpdateTrigger(trigger platform.MultiplayerServerTrigger) Response {
	a.CalledUpdateTrigger = true

	return Response{
		HTTPResponse: &resty.Response{
			RawResponse: &http.Response{
				StatusCode: http.StatusOK,
			},
		},
	}
}

func (a *Spy) DeleteTrigger(id int) Response {
	a.CalledDeleteTrigger = true

	return Response{
		HTTPResponse: &resty.Response{
			RawResponse: &http.Response{
				StatusCode: http.StatusNoContent,
			},
		},
	}
}
