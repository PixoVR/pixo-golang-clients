package auth

type RawToken struct {
	Authorized    bool    `json:"authorized"`
	UserID        int     `json:"userId"`
	OrgType       string  `json:"orgType"`
	OrgID         int     `json:"orgId"`
	FirstName     string  `json:"given_name"`
	LastName      string  `json:"family_name"`
	Email         string  `json:"email"`
	Role          string  `json:"role"`
	Audience      string  `json:"aud"`
	Expiration    int64   `json:"exp"`
	IAT           float64 `json:"iat"`
	Issuer        string  `json:"iss"`
	Sub           string  `json:"sub"`
	JTI           string  `json:"jti"`
	EmailVerified bool    `json:"email_verified"`
	Hd            string  `json:"hd"`
}
