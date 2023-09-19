package auth

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	jwt2 "github.com/go-jose/go-jose/v3/jwt"
	"github.com/golang-jwt/jwt/v4"
	"net/http"
	"os"
	"strings"
)

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
	accessToken := r.Header.Get(SecretKeyHeader)
	if accessToken != "" {
		return accessToken
	}

	authToken := r.Header.Get("Authorization")
	if len(strings.Split(authToken, " ")) == 2 {
		return strings.Split(authToken, " ")[1]
	}

	return ""
}

func extractClaim(claims map[string]interface{}, key string) interface{} {
	if value, ok := claims[key]; ok {
		return value
	}
	return nil
}

func ParseJWT(tokenString string) (RawToken, error) {

	var claims map[string]interface{}

	var token *jwt2.JSONWebToken
	token, err := jwt2.ParseSigned(tokenString)
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
