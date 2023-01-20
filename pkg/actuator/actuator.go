package actuator

import (
	"database/sql"
	"serverless-dbapi/pkg/managercenter"
	"serverless-dbapi/pkg/server"
	"strconv"
	"time"
)

// actuactor config
type Config struct {
	Serer    ServerConfig   `yaml:"server"`
	Database DatabaseConfig `yaml:"database"`
}

// http server config
type ServerConfig struct {
	Port int `yaml:"port"` // port
}

// database config
type DatabaseConfig struct {
	DriverName   string        `yaml:"driver-name"`
	Url          string        `yaml:"url"`            //database url
	MaxIdleCount int           `yaml:"max-idle-count"` // zero means defaultMaxIdleConns; negative means 0
	MaxOpen      int           `yaml:"max-open"`       // <= 0 means unlimited
	MaxLifetime  time.Duration `yaml:"max-lifetime"`   // maximum amount of time a connection may be reused
	MaxIdleTime  time.Duration `yaml:"max-idle-time"`  // maximum amount of time a connection may be idle before being closed
}

// for execing db api
type Actuator struct {
	cfg                 Config
	ser                 server.Server
	dbConn              *sql.DB
	managerCenterServer managercenter.ManagerCenterServer
}

func New(cfg Config) (*Actuator, error) {
	// open connect
	db, err := sql.Open(cfg.Database.DriverName, cfg.Database.Url)
	if err != nil {
		return nil, err
	}
	db.SetMaxIdleConns(cfg.Database.MaxIdleCount)
	db.SetMaxOpenConns(cfg.Database.MaxOpen)
	db.SetConnMaxLifetime(cfg.Database.MaxLifetime)
	db.SetConnMaxIdleTime(cfg.Database.MaxIdleTime)

	return &Actuator{
		dbConn: db,
		cfg:    cfg,
	}, nil
}

func (a *Actuator) SetManagerCenterServer(server managercenter.ManagerCenterServer) *Actuator {
	a.managerCenterServer = server
	return a
}

func (a *Actuator) Run() error {
	// database init
	err := a.dbConn.Ping()
	if err != nil {
		return err
	}

	// http server init
	handle := NewHandle(a.dbConn, a.managerCenterServer)
	a.ser = server.NewActuatorServer(handle.Handler)
	a.ser.Run(":" + strconv.Itoa(a.cfg.Serer.Port))

	return nil
}
