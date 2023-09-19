package auth

import (
	"context"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-jose/go-jose/v3"
	"github.com/golang-jwt/jwt/v4"
	"github.com/rs/zerolog/log"
	"net/http"
	"os"
	"strings"
)

const (
	CustomContextKey = "CUSTOM_CONTEXT"
	GinContextKey    = "GIN_CONTEXT"
	EnforcerKey      = "ENFORCER"
	UserKey          = "USER"
	ErrOrgAccess     = "unable to get organizations for org access"
)

type CustomContext struct {
	Service *services.Service
}

func TokenAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		err := TokenValid(c.Request)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "unauthorized",
			})
			return
		}

		user, err := GetCurrentUser(c)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "unauthorized",
			})
			return
		}
		c.Set(UserKey, user)

		c.Next()
	}
}

func TokenValid(r *http.Request) error {
	tokenString := ExtractSecretKey(r)
	_, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("SECRET_KEY")), nil
	})
	if err != nil {
		return err
	}
	return nil
}

func ExtractSecretKey(r *http.Request) string {
	accessToken := r.Header.Get("x-access-token")
	if accessToken != "" {
		return accessToken
	}

	authToken := r.Header.Get("Authorization")
	if len(strings.Split(authToken, " ")) == 2 {
		return strings.Split(authToken, " ")[1]
	}

	return ""
}

func GetCurrentUser(c *gin.Context) (*models.User, error) {
	context := config.GetContext(c.Request.Context())

	tokenString := ExtractSecretKey(c.Request)

	if tokenString == "" {
		log.Debug().Msg("token not found")
		return nil, errors.New("token not found")
	}

	rawToken, err := ParseJWT(tokenString)
	if err != nil {
		log.Debug().Err(err).Msg("error parsing JWT")
		return nil, err
	}

	user, err := context.Service.UserService.FindByID(rawToken.UserID)
	if err != nil {
		log.Debug().Msg("user not found")
		return nil, err
	}
	user.OrgAccess = GetOrgAccess(context.Service.OrgService, user.Org)

	return user, nil
}

func GetOrgAccess(orgService services.OrgService, userOrg *models.Org) []*models.Org {
	if userOrg == nil {
		return nil
	}

	if userOrg.Type == "platform" {
		allOrgs, err := orgService.GetAll()
		if err != nil {
			log.Debug().Msg(ErrOrgAccess)
			return nil
		}
		return allOrgs
	}

	orgsInTree, err := orgService.GetOrgsInTree(userOrg)
	if err != nil {
		log.Debug().Msg(ErrOrgAccess)
		return nil
	}
	return orgsInTree
}

func GetContext(ctx context.Context) *CustomContext {
	customContext, ok := ctx.Value(CustomContextKey).(*CustomContext)
	if !ok {
		return nil
	}

	return customContext
}

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

func extractClaim(claims map[string]interface{}, key string) interface{} {
	if value, ok := claims[key]; ok {
		return value
	}
	return nil
}

func ParseJWT(tokenString string) (RawToken, error) {

	var claims map[string]interface{}

	token, err := jose.ParseSigned(tokenString)
	if err != nil {
		return RawToken{}, errors.New("no token found")
	}

	_ = token.UnsafeClaimsWithoutVerification(&claims)

	var rawToken RawToken

	if userID, ok := extractClaim(claims, "userId").(float64); ok {
		rawToken.UserID = int(userID)
	}

	if authorized, ok := extractClaim(claims, "authorized").(bool); ok {
		rawToken.Authorized = authorized
	}

	if audience, ok := extractClaim(claims, "aud").(string); ok {
		rawToken.Audience = audience
	}

	if expiration, ok := extractClaim(claims, "exp").(float64); ok {
		rawToken.Expiration = int64(expiration)
	}

	if iat, ok := extractClaim(claims, "iat").(float64); ok {
		rawToken.IAT = iat
	}

	if issuer, ok := extractClaim(claims, "iss").(string); ok {
		rawToken.Issuer = issuer
	}

	if sub, ok := extractClaim(claims, "sub").(string); ok {
		rawToken.Sub = sub
	}

	if jti, ok := extractClaim(claims, "jti").(string); ok {
		rawToken.JTI = jti
	}

	if firstName, ok := extractClaim(claims, "given_name").(string); ok {
		rawToken.FirstName = firstName
	}

	if lastName, ok := extractClaim(claims, "family_name").(string); ok {
		rawToken.LastName = lastName
	}

	if email, ok := extractClaim(claims, "email").(string); ok {
		rawToken.Email = email
	}

	if orgID, ok := extractClaim(claims, "orgId").(int); ok {
		rawToken.OrgID = orgID
	}

	if orgType, ok := extractClaim(claims, "orgType").(string); ok {
		rawToken.OrgType = orgType
	}

	if role, ok := extractClaim(claims, "role").(string); ok {
		rawToken.Role = role
	}

	if emailVerified, ok := extractClaim(claims, "email_verified").(bool); ok {
		rawToken.EmailVerified = emailVerified
	}

	if hd, ok := extractClaim(claims, "hd").(string); ok {
		rawToken.Hd = hd
	}

	return rawToken, nil
}
