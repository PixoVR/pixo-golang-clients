package headset

// DoubleLegacyEventRequest struct
type DoubleLegacyEventRequest struct {
	Type      string  `json:"eventtype,omitempty"`
	SessionID int     `json:"sessionID,omitempty"`
	UserID    int     `json:"userid,omitempty"`
	ModuleID  int     `json:"moduleID,omitempty"`
	Payload   Payload `json:"jsonData,omitempty"`
}
