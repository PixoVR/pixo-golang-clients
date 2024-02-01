package config

type inMemoryManagerImpl struct {
	baseManagerImpl
	config Config
}

func NewInMemoryManager() Manager {
	return &inMemoryManagerImpl{}
}

func (m *inMemoryManagerImpl) SetConfig(config Config) {
	m.config = config
	m.activeEnv = config.ActiveEnv()
}

func (m *inMemoryManagerImpl) GetConfig() Config {
	return m.config
}

func (m *inMemoryManagerImpl) ReadConfigFile(configFile string) error {
	m.SetConfigFile(configFile)
	return nil
}

func (m *inMemoryManagerImpl) ConfigFile() string {
	return m.configFile
}

func (m *inMemoryManagerImpl) SetConfigValue(key, value string) error {
	return nil
}

func (m *inMemoryManagerImpl) GetConfigValue(key string) (string, bool) {
	env := m.GetActiveEnv()
	return env.Get(key)
}

func (m *inMemoryManagerImpl) GetActiveEnv() Env {
	config := m.GetConfig()

	return config.ActiveEnv()
}

func (m *inMemoryManagerImpl) SetActiveEnv(env Env) {
	config := m.GetConfig()

	m.config.Region = env.Region
	m.config.Lifecycle = env.Lifecycle
	m.activeEnv = config.Envs[env.Name()]
}

func (m *inMemoryManagerImpl) GetOrAsk(key string) string {
	return ""
}
