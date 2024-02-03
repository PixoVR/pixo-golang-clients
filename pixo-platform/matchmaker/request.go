package matchmaker

import (
	"strconv"
)

type MatchRequest struct {
	ModuleID      int    `json:"moduleId"`
	ServerVersion string `json:"serverVersion"`
}

func (m MatchRequest) IsValid() bool {
	return m.ModuleID != 0 && m.ServerVersion != ""
}

type MatchResponse struct {
	Error        bool         `json:"error,omitempty"`
	Message      string       `json:"message,omitempty"`
	MatchDetails MatchDetails `json:"matchDetails,omitempty"`
}

type MatchDetails struct {
	IP             string `json:"IPAddress,omitempty"`
	Port           string `json:"Port,omitempty"`
	SessionName    string `json:"SessionName,omitempty"`
	SessionID      string `json:"SessionID,omitempty"`
	MapName        string `json:"MapName,omitempty"`
	OwningUserName string `json:"OwningUserName,omitempty"`
	ModuleVersion  string `json:"ServerVersion,omitempty"`
	ModuleID       int    `json:"ModuleID,omitempty"`
	OrgID          int    `json:"OrgID,omitempty"`
}

func (m MatchResponse) IsValid() bool {
	if m.Error {
		return false
	}

	if m.MatchDetails.IP == "" || m.MatchDetails.Port == "" {
		return false
	}

	if _, err := strconv.Atoi(m.MatchDetails.Port); err != nil {
		return false
	}

	return true
}
