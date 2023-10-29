package matchmaker

import (
	"strconv"
)

type MatchRequest struct {
	OrgID    int `json:"orgId"`
	ModuleID int `json:"moduleId"`
}

func (m MatchRequest) IsValid() bool {
	return m.OrgID != 0 && m.ModuleID != 0
}

type MatchResponse struct {
	Error        bool         `json:"error"`
	Message      string       `json:"message"`
	MatchDetails MatchDetails `json:"matchDetails"`
}

type MatchDetails struct {
	IPAddress      string `json:"IPAddress"`
	Port           string `json:"Port"`
	SessionName    string `json:"SessionName"`
	SessionID      string `json:"SessionID"`
	MapName        string `json:"MapName"`
	OwningUserName string `json:"OwningUserName"`
	OrgID          int    `json:"OrgId"`
	ModuleID       int    `json:"ModuleId"`
	ModuleVersion  string `json:"ModuleVersion"`
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
