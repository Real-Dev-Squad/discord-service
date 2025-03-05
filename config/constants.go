package config

import (
	"fmt"

	"github.com/sirupsen/logrus"
)

type Environment string

const (
	Development Environment = "development"
	Staging     Environment = "staging"
	Production  Environment = "production"
)

func (e Environment) Validate() Environment {
	if e != Development && e != Staging && e != Production {
		logrus.Warn("Invalid Environment. Defaulting to development.")
		e = Development
	}
	return e
}

type URLs struct {
	RDS_BASE_API_URL      string
	VERIFICATION_SITE_URL string
}

const domain = "https://%s.realdevsquad.com"

var EnvironmentURLs = map[Environment]URLs{
	Development: {
		RDS_BASE_API_URL:      "http://localhost:3000",
		VERIFICATION_SITE_URL: "http://localhost:3443",
	},
	Staging: {
		RDS_BASE_API_URL:      fmt.Sprintf(domain, "staging-api"),
		VERIFICATION_SITE_URL: fmt.Sprintf(domain, "staging-my"),
	},
	Production: {
		RDS_BASE_API_URL:      fmt.Sprintf(domain, "api"),
		VERIFICATION_SITE_URL: fmt.Sprintf(domain, "my"),
	},
}
