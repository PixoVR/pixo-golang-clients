package multiplayer_allocator

import "github.com/go-resty/resty/v2"

type AllocationRequest struct {
	ModuleID           int    `json:"module_id"`
	OrgID              int    `json:"org_id"`
	ImageRegistry      string `json:"image_registry"`
	Engine             string `json:"engine"`
	BackfillID         string `json:"backfill_id"`
	AllocateGameServer bool   `json:"allocate_game_server"`
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
