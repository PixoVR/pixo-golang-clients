package multiplayer_allocator

import "os"

func getURL() string {
	apiURL, ok := os.LookupEnv("ALLOCATOR_API_URL")
	if !ok {
		return "https://multi-central1.allocator.multiplayer.pixovr.com"
	}

	return apiURL
}
