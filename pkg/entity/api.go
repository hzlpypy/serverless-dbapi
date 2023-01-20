package entity

import "time"

// api type
const (
	DATA_API_TYPE = 1
)

// datasource info
type DatabaseConfig struct {
	Id           string        `yaml:"id"`
	Name         string        `yaml:"name"`
	DriverName   string        `yaml:"driver-name"`
	Url          string        `yaml:"url"`            //database url
	MaxIdleCount int           `yaml:"max-idle-count"` // zero means defaultMaxIdleConns; negative means 0
	MaxOpen      int           `yaml:"max-open"`       // <= 0 means unlimited
	MaxLifetime  time.Duration `yaml:"max-lifetime"`   // maximum amount of time a connection may be reused
	MaxIdleTime  time.Duration `yaml:"max-idle-time"`  // maximum amount of time a connection may be idle before being closed
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
	return DATASOURCE_PREFIX + d.Id
}

// api group info
type ApiGroupConfig struct {
	Id   string
	Name string
}

func (d *ApiGroupConfig) SetId(id string) {
	d.Id = id
}

func (d *ApiGroupConfig) GetId() string {
	return d.Id
}

func (d *ApiGroupConfig) GetPrefixId() string {
	return API_GROUP_PREFIX + d.Id
}

// api config info
type ApiConfig struct {
	Id           string
	ApiGroupId   string
	ApiType      int
	Sql          string
	ParamKey     []string
	DataSourceId string
}

func (d *ApiConfig) SetId(id string) {
	d.Id = id
}

func (d *ApiConfig) GetId() string {
	return d.Id
}

// for quickly select
func (d *ApiConfig) GetPrefixId() string {
	return API_GROUP_PREFIX + d.ApiGroupId + API_PREFIX + d.Id
}

type IdCommon interface {
	GetId() string
	SetId(id string)
	GetPrefixId() string
}
