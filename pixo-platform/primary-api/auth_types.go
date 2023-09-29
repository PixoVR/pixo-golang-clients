package primary_api

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
	ID         int    `json:"id"`
	OrgID      int    `json:"orgId"`
	First      string `json:"first"`
	Last       string `json:"last"`
	Email      string `json:"email"`
	ExternalID string `json:"externalId"`
}
