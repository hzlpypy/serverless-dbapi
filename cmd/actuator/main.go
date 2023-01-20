package main

import (
	"os"
	"serverless-dbapi/pkg/actuator"
	"serverless-dbapi/pkg/cfg"
	"serverless-dbapi/pkg/managercenter"
	"serverless-dbapi/pkg/mode"

	_ "github.com/go-sql-driver/mysql"
	"gopkg.in/yaml.v3"
)

func main() {
	// load config file
	file, err := os.ReadFile("./actuator.yaml")
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

	// new actuator
	actuator, err := actuator.New(config.Actuator)
	actuator.SetManagerCenterServer(managercenter.NewManagerCenterServer())

	if err != nil {
		panic(err)
	}
	err = actuator.Run()
	if err != nil {
		panic(err)
	}
}
