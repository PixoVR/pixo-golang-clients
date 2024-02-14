package primary_api

import "time"

type LoginRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type LegacyAuthResponse struct {
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
	FirstName  string `json:"firstName,omitempty"`
	LastName   string `json:"lastName,omitempty"`
	Email      string `json:"email,omitempty"`
	ExternalID string `json:"externalId,omitempty"`

	CreatedAt time.Time `json:"createdAt,omitempty"`
	UpdatedAt time.Time `json:"updatedAt,omitempty"`
}

type APIKey struct {
	ID     int    `json:"id,omitempty"`
	Key    string `json:"key,omitempty"`
	UserID int    `json:"userId,omitempty"`
	User   *User  `json:"user,omitempty"`

	CreatedAt time.Time `json:"createdAt,omitempty"`
	UpdatedAt time.Time `json:"updatedAt,omitempty"`
}
