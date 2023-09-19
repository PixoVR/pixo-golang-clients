package auth

import "os"

func IsValidSecretKey(input string) bool {
	key := os.Getenv("SECRET_KEY")
	return input != "" && input == key
}
