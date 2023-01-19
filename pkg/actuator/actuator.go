package actuator

import (
	"database/sql"
	"serverless-db/pkg/server"
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
	cfg    Config
	ser    server.Server
	dbConn *sql.DB
}

func New(cfg Config) (Actuator, error) {
	return Actuator{
		cfg: cfg,
	}, nil
}

func (a *Actuator) Run() error {
	// database init
	db, err := sql.Open(a.cfg.Database.DriverName, a.cfg.Database.Url)
	if err != nil {
		return err
	}
	db.SetMaxIdleConns(a.cfg.Database.MaxIdleCount)
	db.SetMaxOpenConns(a.cfg.Database.MaxOpen)
	db.SetConnMaxLifetime(a.cfg.Database.MaxLifetime)
	db.SetConnMaxIdleTime(a.cfg.Database.MaxIdleTime)
	err = db.Ping()
	if err != nil {
		return err
	}
	a.dbConn = db

	// http server init
	handle := NewHandle(a.dbConn, &server.MockManagerCenterServer{})
	a.ser = server.NewActuatorServer(handle.Handler)
	a.ser.Run(":" + strconv.Itoa(a.cfg.Serer.Port))

	return nil
}
