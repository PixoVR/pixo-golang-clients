package urlfinder

import (
	"fmt"
)

type ServiceConfig struct {
	ServiceName string
	Service     string
	Port        int
	Tenant      string
	Region      string
	Lifecycle   string
	Namespace   string
	InternalDNS bool
}

type ClientConfig struct {
	APIKey    string
	Token     string
	Internal  bool
	Lifecycle string
	Namespace string
	Region    string
}

func (s ServiceConfig) FormatURL() string {

	if s.Service == "" {
		s.Service = DefaultService
	}

	if s.Lifecycle == "local" {
		if s.Service == "matchmaking" {
			return fmt.Sprintf("ws://localhost:%d", s.Port)
		}

		if s.Port == 0 {
			s.Port = 8000
		}

		return fmt.Sprintf("http://localhost:%d/%s", s.Port, s.Service)
	}

	if s.Tenant == "" {
		s.Tenant = DefaultTenant
	}

	if s.Region == "" {
		s.Region = DefaultRegion
	}

	if s.InternalDNS {
		return fmt.Sprintf("http://%s-%s.%s.svc", s.Namespace, s.ServiceName, s.Namespace)
	}

	var prefix string
	switch s.Region {
	case "na":
		if s.Tenant == "multiplayer" {
			prefix = "multi-central1"
		}
	case "saudi":
		if s.Lifecycle == "prod" || s.Lifecycle == "" {
			if s.Tenant == "apex" && s.Service == "v2" {
				return fmt.Sprintf("https://apisa.pixovr.com/v2")
			} else if s.Service == "api" {
				return fmt.Sprintf("https://apisa.pixovr.com")
			}
		}

		if s.Tenant == "multiplayer" {
			prefix = "multi-saudi"
		} else {
			prefix = "saudi"
		}
	}

	if prefix != "" {
		s.Tenant = fmt.Sprintf("%s.%s", prefix, s.Tenant)
	}

	if s.Service == "api" {
		return fmt.Sprintf("https://%s.%s.%s", s.Service, s.Tenant, s.GetBaseDomain())
	}

	protocol := "https"
	if s.Service == "matchmaking" {
		protocol = "wss"
	}

	return fmt.Sprintf("%s://%s.%s/%s", protocol, s.Tenant, s.GetBaseDomain(), s.Service)
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
