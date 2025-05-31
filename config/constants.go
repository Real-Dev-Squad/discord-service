package config

import (
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

var EnvironmentURLs = map[Environment]URLs{
	Development: {
		RDS_BASE_API_URL:      AppConfig.RDS_BASE_API_URL,
		VERIFICATION_SITE_URL: AppConfig.VERIFICATION_SITE_URL,
	},
	Staging: {
		RDS_BASE_API_URL:      AppConfig.RDS_BASE_API_URL,
		VERIFICATION_SITE_URL: AppConfig.VERIFICATION_SITE_URL,
	},
	Production: {
		RDS_BASE_API_URL:      AppConfig.RDS_BASE_API_URL,
		VERIFICATION_SITE_URL: AppConfig.VERIFICATION_SITE_URL,
	},
}
