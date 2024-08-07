package urlfinder

import (
	"fmt"
	"strings"
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

func (s ServiceConfig) FormatURL(protocolInput ...string) string {

	protocol := "https"
	if len(protocolInput) > 0 {
		protocol = protocolInput[0]
		if strings.HasPrefix(protocol, "ws") {
			protocol = "wss"
		}
	}

	if s.Service == "" {
		s.Service = DefaultService
	}

	if s.Lifecycle == "local" {
		if protocol == "wss" {
			protocol = "ws"
		} else {
			protocol = "http"
		}

		if s.Port == 0 {
			s.Port = 8000
		}

		if s.Service == "api" || s.Service == "modules" {
			return fmt.Sprintf("http://localhost:%d", s.Port)
		}

		return fmt.Sprintf("%s://localhost:%d/%s", protocol, s.Port, s.Service)
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

	if prefix != "" {
		if s.Service == "api" {
			prefix = fmt.Sprintf("%s.%s", prefix, s.Service)
		} else {
			prefix = fmt.Sprintf("%s.%s", prefix, s.Tenant)
		}
	} else {
		if s.Service == "api" || s.Service == "modules" {
			prefix = s.Service
		}
	}

	if s.Service == "api" || s.Service == "modules" {
		return fmt.Sprintf("https://%s.%s.%s", prefix, s.Tenant, s.GetBaseDomain())
	}

	if prefix == "" {
		prefix = s.Tenant
	}

	return fmt.Sprintf("%s://%s.%s/%s", protocol, prefix, s.GetBaseDomain(), s.Service)
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
