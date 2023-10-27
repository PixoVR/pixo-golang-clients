package multiplayer_allocator

import (
	primary_api "github.com/PixoVR/pixo-golang-clients/pixo-platform/primary-api"
	"github.com/go-resty/resty/v2"
)

type FleetRequest struct {
	ServerVersion   primary_api.MultiplayerServerVersion `json:"serverVersion,omitempty"`
	StandbyReplicas int                                  `json:"standbyReplicas,omitempty"`
}

type Response struct {
	HTTPResponse *resty.Response `json:"http_response"`
	Error        error           `json:"error"`
}

type AllocationRequest struct {
	ModuleID           int    `json:"module_id,omitempty"`
	OrgID              int    `json:"org_id,omitempty"`
	ImageRegistry      string `json:"image_registry,omitempty"`
	Engine             string `json:"engine,omitempty"`
	BackfillID         string `json:"backfill_id,omitempty"`
	AllocateGameServer bool   `json:"allocate_game_server,omitempty"`
	ServerVersion      string `json:"server_version,omitempty"`
}

type AllocationResponse struct {
	HTTPResponse *resty.Response `json:"http_response"`
	Results      GameServer      `json:"results"`
	Error        error           `json:"error"`
}

type GameServer struct {
	ResourceName   string `json:"resource_name"`
	Address        string `json:"ipaddress"`
	Port           string `json:"port"`
	SessionName    string `json:"session_name"`
	SessionID      string `json:"session_id"`
	OwningUserName string `json:"owning_user_name"`
	OrgID          int    `json:"org_id"`
	ModuleID       int    `json:"module_id"`
	ModuleVersion  string `json:"module_version"`
	MapName        string `json:"map_name"`
	State          string `json:"state"`
	NumPlaying     int    `json:"num_playing"`
	NumBackfill    int    `json:"num_backfill"`
	Capacity       int    `json:"capacity"`
	Logs           string `json:"logs"`
}
