package env

import (
	"fmt"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"
	"path/filepath"
	"runtime"
)

func GetProjectRoot() (root string) {
	_, b, _, _ := runtime.Caller(0)
	root = filepath.Join(filepath.Dir(b), "..")
	return root
}

func SourceProjectEnv() {
	filename := fmt.Sprintf("%s/.env", GetProjectRoot())
	if err := godotenv.Load(filename); err != nil {
		log.Warn().
			Err(err).
			Str("file", filename).
			Msg("file not found")
	}
}
