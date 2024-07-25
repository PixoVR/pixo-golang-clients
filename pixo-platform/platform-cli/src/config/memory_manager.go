package config

var _ Manager = (*InMemoryConfigManager)(nil)

type InMemoryConfigManager struct {
	config *Config
}

func NewInMemoryConfigManager() *InMemoryConfigManager {
	return &InMemoryConfigManager{config: &Config{}}
}

func (i *InMemoryConfigManager) GetConfig() *Config {
	return i.config
}

func (i *InMemoryConfigManager) SetConfig(config Config) {
	i.config = &config
}

func (i *InMemoryConfigManager) Clear() {
	i.config = &Config{}
}
