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
	Key       string
	Token     string
	Internal  bool
	Lifecycle string
	Namespace string
	Region    string
}

func (s ServiceConfig) FormatURL() string {

	if s.Lifecycle == "local" {
		if s.Service == "matchmaking" {
			return fmt.Sprintf("ws://localhost:%d", s.Port)
		}

		if s.Port == 0 {
			s.Port = 8000
		}

		if s.Service == "" {
			s.Service = "v2"
		}

		return fmt.Sprintf("http://localhost:%d/%s", s.Port, s.Service)
	}

	if s.Service == "" {
		s.Service = DefaultService
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
		if s.Tenant == "multiplayer" {
			prefix = "multi-saudi"
		} else {
			prefix = "saudi"
		}
	}

	protocol := "https"
	if s.Service == "matchmaking" {
		protocol = "wss"
	}

	if prefix != "" {
		s.Tenant = fmt.Sprintf("%s.%s", prefix, s.Tenant)
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
