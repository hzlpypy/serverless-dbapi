package cfg

import "serverless-dbapi/pkg/actuator"

// all config
type Config struct {
	Mode     string          `yaml:"mode"`
	Actuator actuator.Config `yaml:"actuator"`
}
