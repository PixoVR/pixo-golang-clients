package multiplayer_allocator

import (
	platform "github.com/PixoVR/pixo-golang-clients/pixo-platform/primary-api"
	"github.com/go-resty/resty/v2"
	"net/http"
)

type BuildTriggerSpy struct {
	CalledRegisterTrigger bool
	CalledUpdateTrigger   bool
	CalledDeleteTrigger   bool
}

func NewBuildTriggerSpy() *BuildTriggerSpy {
	return &BuildTriggerSpy{
		CalledRegisterTrigger: false,
		CalledUpdateTrigger:   false,
		CalledDeleteTrigger:   false,
	}
}

func (b *BuildTriggerSpy) RegisterTrigger(trigger platform.MultiplayerServerTrigger) (*resty.Response, error) {
	b.CalledRegisterTrigger = true

	return &resty.Response{
		RawResponse: &http.Response{
			StatusCode: http.StatusCreated,
		},
	}, nil
}

func (b *BuildTriggerSpy) UpdateTrigger(trigger platform.MultiplayerServerTrigger) (*resty.Response, error) {
	b.CalledUpdateTrigger = true

	return &resty.Response{
		RawResponse: &http.Response{
			StatusCode: http.StatusOK,
		},
	}, nil
}

func (b *BuildTriggerSpy) DeleteTrigger(id int) (*resty.Response, error) {
	b.CalledDeleteTrigger = true

	return &resty.Response{
		RawResponse: &http.Response{
			StatusCode: http.StatusNoContent,
		},
	}, nil
}
