package abstract_client

import "os"

func getAPIURL() string {
	apiURL, ok := os.LookupEnv("API_URL")
	if !ok {
		return "https://api.apex.dev.pixovr.com"
	}

	return apiURL
}
