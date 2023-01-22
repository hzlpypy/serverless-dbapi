package main

import (
	"errors"
	"os"
	"serverless-dbapi/cmd/mode"
	"serverless-dbapi/pkg/actuator"
	"serverless-dbapi/pkg/cfg"
	"serverless-dbapi/pkg/managercenter"

	_ "github.com/go-sql-driver/mysql"
	"gopkg.in/yaml.v3"
)

func main() {
	// load config file
	file, err := os.ReadFile("../config.yaml")
	if err != nil {
		panic(err)
	}
	config := cfg.Config{}
	err = yaml.Unmarshal(file, &config)
	if err != nil {
		panic(err)
	}

	// mark server mode: mock、standalone、cluster
	mode.MODE = config.Mode

	// new manager center
	managerCenterServer, err := newManagerCenterServer(config.ManagerCenter)
	if err != nil {
		panic(err)
	}

	if mode.MOCK != mode.CLUSTER {
		// if the mode is not cluster, get all databases for manager center
		// if the mode is cluster, just get from config
		config.Actuator.Databases, err = managerCenterServer.GetDataBases()
		if err != nil {
			panic(err)
		}
	}

	// new actuator
	actuator, err := actuator.New(config.Actuator)
	if err != nil {
		panic(err)
	}
	actuator.SetManagerCenterServer(managerCenterServer)

	err = actuator.Run(config.Serer.Port)
	if err != nil {
		panic(err)
	}
}

func newManagerCenterServer(config *cfg.ManagerCenterConfig) (managercenter.ManagerCenterServer, error) {
	if mode.MODE == mode.MOCK {
		return managercenter.NewMockManagerCenterServer()
	}
	if mode.MODE == mode.STANDALONE {
		return managercenter.NewMemoryManagerCenterServer(config)
	}
	if mode.MODE == mode.CLUSTER {
		return managercenter.NewHttpManagerCenterServer()
	}
	return nil, errors.New("manager center server not found")
}
