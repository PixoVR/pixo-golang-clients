package config

import (
	"fmt"
	"io"
	"os"
	"strings"
)

type fileManagerImpl struct {
	configFile string
	reader     io.Reader
	writer     io.Writer
}

func NewFileManager(cfgFile string) Manager {
	var m fileManagerImpl

	if cfgFile != "" {
		_ = m.SetConfigFile(cfgFile)
	}

	return &m
}

func (f *fileManagerImpl) Lifecycle() string {
	if lifecycle, ok := os.LookupEnv("PIXO_LIFECYCLE"); ok {
		return lifecycle
	}

	return f.GetActiveEnv().Lifecycle
}

func (f *fileManagerImpl) Region() string {
	if region, ok := os.LookupEnv("PIXO_REGION"); ok {
		return region
	}

	return f.GetActiveEnv().Region
}

func (f *fileManagerImpl) formatEnvVar(key string) string {
	return fmt.Sprintf("PIXO_%s", strings.ReplaceAll(strings.ToUpper(key), "-", "_"))
}
