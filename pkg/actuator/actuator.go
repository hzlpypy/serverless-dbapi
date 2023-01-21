package actuator

import (
	"database/sql"
	"serverless-dbapi/pkg/cfg"
	"serverless-dbapi/pkg/managercenter"
	"serverless-dbapi/pkg/server"
	"strconv"
)

// for execing db api
type Actuator struct {
	cfg                 *cfg.ActuactorConfig
	ser                 server.Server
	dbConns             map[string]*sql.DB
	managerCenterServer managercenter.ManagerCenterServer
}

func New(cfg *cfg.ActuactorConfig) (*Actuator, error) {
	// open connect
	dbMap := make(map[string]*sql.DB)
	if len(cfg.Databases) > 0 {
		for _, value := range cfg.Databases {
			db, err := sql.Open(value.DriverName, value.Url)
			if err != nil {
				return nil, err
			}
			db.SetMaxIdleConns(value.MaxIdleCount)
			db.SetMaxOpenConns(value.MaxOpen)
			db.SetConnMaxLifetime(value.MaxLifetime)
			db.SetConnMaxIdleTime(value.MaxIdleTime)
			dbMap[value.Id] = db
		}
	}

	return &Actuator{
		dbConns: dbMap,
		cfg:     cfg,
	}, nil
}

func (a *Actuator) SetManagerCenterServer(server managercenter.ManagerCenterServer) *Actuator {
	a.managerCenterServer = server
	return a
}

func (a *Actuator) Run(port int) error {
	// database init
	if len(a.dbConns) > 0 {
		for _, value := range a.dbConns {
			err := value.Ping()
			if err != nil {
				return err
			}
		}
	}

	// http server init
	handle := NewHandle(a.dbConns, a.managerCenterServer)
	a.ser = server.NewActuatorServer(handle.Handler)
	a.ser.Run(":" + strconv.Itoa(port))

	return nil
}
