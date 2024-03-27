package config

import (
	"fmt"
	"github.com/PixoVR/pixo-golang-clients/pixo-platform/platform-cli/pkg/forms"
	"github.com/PixoVR/pixo-golang-clients/pixo-platform/platform-cli/pkg/forms/basic"
	"io"
	"os"
	"strings"
)

type fileManagerImpl struct {
	configFile  string
	reader      io.Reader
	writer      io.Writer
	formHandler forms.FormHandler
}

func NewFileManager(cfgFile string, formHandlers ...forms.FormHandler) Manager {
	var formHandler forms.FormHandler
	if len(formHandlers) > 0 {
		formHandler = formHandlers[0]
	} else {
		formHandler = basic.NewFormHandler(os.Stdin, os.Stdout)
	}

	m := &fileManagerImpl{
		formHandler: formHandler,
	}

	if cfgFile != "" {
		_ = m.SetConfigFile(cfgFile)
	}

	return m
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
