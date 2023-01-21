package entity

import (
	"serverless-dbapi/pkg/tool"
	"time"
)

// api type
const (
	DATA_API_TYPE = 1
)

// datasource info
type DatabaseConfig struct {
	Id           string        `yaml:"id" json:"id,omitempty"`
	Name         string        `yaml:"name" json:"name,omitempty" validate:"required"`
	DriverName   string        `yaml:"driver-name" json:"driverName,omitempty" validate:"required"`
	Url          string        `yaml:"url" json:"url,omitempty" validate:"required"` //database url
	MaxIdleCount int           `yaml:"max-idle-count" json:"maxIdleCount"`           // zero means defaultMaxIdleConns; negative means 0
	MaxOpen      int           `yaml:"max-open" json:"maxOpen"`                      // <= 0 means unlimited
	MaxLifetime  time.Duration `yaml:"max-lifetime" json:"maxLifetime"`              // maximum amount of time a connection may be reused
	MaxIdleTime  time.Duration `yaml:"max-idle-time" json:"maxIdletime"`             // maximum amount of time a connection may be idle before being closed
}

const (
	DATASOURCE_PREFIX = "/datasource/"
	API_GROUP_PREFIX  = "/api-group/"
	API_PREFIX        = "/api/"
)

func (d *DatabaseConfig) SetId(id string) {
	d.Id = id
}

func (d *DatabaseConfig) GetId() string {
	return d.Id
}

func (d *DatabaseConfig) GetPrefixId() string {
	return tool.StringBuilder(DATASOURCE_PREFIX, d.Id)
}

// api group info
type ApiGroupConfig struct {
	Id   string `json:"id,omitempty"`
	Name string `json:"name,omitempty" validate:"required"`
}

func (d *ApiGroupConfig) SetId(id string) {
	d.Id = id
}

func (d *ApiGroupConfig) GetId() string {
	return d.Id
}

func (d *ApiGroupConfig) GetPrefixId() string {
	return tool.StringBuilder(API_GROUP_PREFIX, d.Id)
}

// api config info
type ApiConfig struct {
	Id           string   `json:"id,omitempty"`
	ApiGroupId   string   `json:"apiGroupId,omitempty" validate:"required"`
	ApiType      int      `json:"apiType,omitempty" validate:"required"`
	Sql          string   `json:"sql,omitempty" validate:"required"`
	ParamKey     []string `json:"paramKey,omitempty" validate:"required"`
	DataSourceId string   `json:"datasourceId,omitempty" validate:"required"`
}

func (d *ApiConfig) SetId(id string) {
	d.Id = id
}

func (d *ApiConfig) GetId() string {
	return d.Id
}

// for quickly select
func (d *ApiConfig) GetPrefixId() string {
	return tool.StringBuilder(API_GROUP_PREFIX, d.ApiGroupId, API_PREFIX, d.Id)
}

type IdCommon interface {
	GetId() string
	SetId(id string)
	GetPrefixId() string
}
