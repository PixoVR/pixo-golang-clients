package multiplayer_allocator

import "os"

func getURL() string {
	apiURL, ok := os.LookupEnv("API_URL")
	if !ok {
		return "https://multi-central1.allocator.multiplayer.dev.pixovr.com"
	}

	return apiURL
}
