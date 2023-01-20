package cfg

import (
	"serverless-dbapi/pkg/entity"
)

// all config
type Config struct {
	Mode     string          `yaml:"mode"`
	Actuator ActuactorConfig `yaml:"actuator"`
}

// actuactor config
type ActuactorConfig struct {
	Serer     ServerConfig            `yaml:"server"`
	Databases []entity.DatabaseConfig `yaml:"databases"`
}

// http server config
type ServerConfig struct {
	Port int `yaml:"port"` // port
}
