package primary_api

import "os"

func getURL() string {
	apiURL, ok := os.LookupEnv("API_URL")
	if !ok {
		return "https://api.apex.pixovr.com"
	}

	return apiURL
}
