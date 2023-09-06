package vertexai

import (
	"context"
	"github.com/rs/zerolog/log"

	"golang.org/x/oauth2"
	auth "golang.org/x/oauth2/google"
)

func getAccessToken() (string, error) {
	var token *oauth2.Token
	ctx := context.Background()

	scopes := []string{
		"https://www.googleapis.com/auth/cloud-platform",
	}

	credentials, err := auth.FindDefaultCredentials(ctx, scopes...)
	if err != nil {
		log.Error().Err(err)
		return "", err
	}

	token, err = credentials.TokenSource.Token()
	if err != nil {
		log.Error().Err(err)
		return "", err
	}

	return token.AccessToken, nil
}
