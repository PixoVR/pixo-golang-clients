package graphql_api

import "os"

func getURL() string {
	apiURL, ok := os.LookupEnv("API_URL")
	if !ok {
		return "https://primary.apex.pixovr.com"
	}

	return apiURL
}
