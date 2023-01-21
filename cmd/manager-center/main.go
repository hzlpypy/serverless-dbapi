package main

import (
	"os"
	"serverless-dbapi/pkg/cfg"
	"serverless-dbapi/pkg/managercenter"

	"gopkg.in/yaml.v2"
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
	server, err := managercenter.New(config.ManagerCenter)
	if err != nil {
		panic(err)
	}

	err = server.Run(config.Serer.Port)
	if err != nil {
		panic(err)
	}
}
