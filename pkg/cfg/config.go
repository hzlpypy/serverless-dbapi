package cfg

import (
	"serverless-dbapi/pkg/entity"
)

// all config
type Config struct {
	Serer         *ServerConfig        `yaml:"server"`
	Mode          string               `yaml:"mode"`
	Actuator      *ActuactorConfig     `yaml:"actuator"`
	ManagerCenter *ManagerCenterConfig `yaml:"manager-center"`
}

// actuactor config
type ActuactorConfig struct {
	Databases []*entity.DatabaseConfig `yaml:"databases"`
}

// http server config
type ServerConfig struct {
	Port int `yaml:"port"` // port
}

// manager center config
type ManagerCenterConfig struct {
	Store *StoreConfig `yaml:"store"`
}

type StoreConfig struct {
	Etcd *EtcdConfig `yaml:"etcd"`
}

type EtcdConfig struct {
	Prefix    string   `yaml:"prefix"`
	Endpoints []string `yaml:"endpoints"`
}
