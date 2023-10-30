package multiplayer_allocator

import (
	platform "github.com/PixoVR/pixo-golang-clients/pixo-platform/primary-api"
	"github.com/go-resty/resty/v2"
)

type FleetRequest struct {
	ServerConfig  platform.MultiplayerServerConfig  `json:"serverConfig,omitempty"`
	ServerVersion platform.MultiplayerServerVersion `json:"serverVersion,omitempty"`
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
	Name           string `json:"resource_name,omitempty"`
	IP             string `json:"ipaddress,omitempty"`
	Port           string `json:"port,omitempty"`
	SessionName    string `json:"session_name,omitempty"`
	SessionID      string `json:"session_id,omitempty"`
	OwningUserName string `json:"owning_user_name,omitempty"`
	OrgID          int    `json:"org_id,omitempty"`
	ModuleID       int    `json:"module_id,omitempty"`
	ModuleVersion  string `json:"module_version,omitempty"`
	MapName        string `json:"map_name,omitempty"`
	State          string `json:"state,omitempty"`
	NumPlaying     int    `json:"num_playing,omitempty"`
	NumBackfill    int    `json:"num_backfill,omitempty"`
	Capacity       int    `json:"capacity,omitempty"`
	Logs           string `json:"logs,omitempty"`
	SidecarLogs    string `json:"sidecar_logs,omitempty"`
}
