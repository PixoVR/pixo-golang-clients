package config

import "io"

type baseManagerImpl struct {
	configFile string
	reader     io.Reader
	writer     io.Writer
	activeEnv  Env
}

func (m *baseManagerImpl) SetReader(r io.Reader) {
	m.reader = r
}

func (m *baseManagerImpl) Reader() io.Reader {
	return m.reader
}

func (m *baseManagerImpl) SetWriter(w io.Writer) {
	m.writer = w
}

func (m *baseManagerImpl) Writer() io.Writer {
	return m.writer
}

func (m *baseManagerImpl) SetConfigFile(configFile string) {
	m.configFile = configFile
}

func (m *baseManagerImpl) Lifecycle() string {
	return m.activeEnv.Lifecycle
}

func (m *baseManagerImpl) Region() string {
	return m.activeEnv.Region
}
