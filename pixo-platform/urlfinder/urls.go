package urlfinder

import (
	"fmt"
)

type ServiceConfig struct {
	Service   string
	Port      int
	Tenant    string
	Region    string
	Lifecycle string
}

func (s ServiceConfig) FormatURL() string {

	if s.Lifecycle == "local" {
		if s.Service == "match" {
			return fmt.Sprintf("ws://localhost:%d", s.Port)
		}

		return fmt.Sprintf("http://localhost:%d", s.Port)
	}

	if s.Service == "" {
		s.Service = DefaultAPIService
	}

	if s.Tenant == "" {
		s.Tenant = DefaultAPITenant
	}

	if s.Region == "" {
		s.Region = DefaultAPIRegion
	}

	var prefix string
	switch s.Region {
	case "na":
		if s.Tenant == "multiplayer" {
			prefix = "multi-central1"
		}
	case "saudi":
		if s.Tenant == "multiplayer" {
			prefix = "multi-saudi"
		} else {
			prefix = "saudi"
		}
	}

	protocol := "https"
	if s.Service == "match" {
		protocol = "wss"
	}

	if prefix != "" {
		s.Service = fmt.Sprintf("%s.%s", prefix, s.Service)
	}

	return fmt.Sprintf("%s://%s.%s.%s", protocol, s.Service, s.Tenant, s.GetBaseDomain())
}

func (s ServiceConfig) GetBaseDomain() string {
	if s.Lifecycle == "local" {
		return "localhost"
	}

	if s.Lifecycle != "" && s.Lifecycle != "prod" {
		return fmt.Sprintf("%s.%s", s.Lifecycle, Domain)
	}

	return Domain
}
