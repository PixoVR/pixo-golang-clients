package primary_api

import "time"

type LoginRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type AuthResponse struct {
	User struct {
		Token string `json:"authToken"`
		Role  string `json:"role"`
	}
}

type User struct {
	ID         int    `json:"id,omitempty"`
	Role       string `json:"role,omitempty"`
	OrgID      int    `json:"orgId,omitempty"`
	Org        Org    `json:"org,omitempty"`
	Username   string `json:"username,omitempty"`
	Password   string `json:"-"`
	FirstName  string `json:"first,omitempty"`
	LastName   string `json:"last,omitempty"`
	Email      string `json:"email,omitempty"`
	ExternalID string `json:"externalId,omitempty"`

	CreatedAt time.Time `json:"createdAt,omitempty"`
	UpdatedAt time.Time `json:"updatedAt,omitempty"`
}
