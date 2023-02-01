package main

import (
	"context"
	"os"
	"serverless-dbapi/pkg/cfg"
	"serverless-dbapi/pkg/managercenter"
	"time"

	"github.com/coreos/etcd/clientv3"
	edclient "github.com/kiraqjx/ed-client"
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

	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   config.Discovery.Etcd.Endpoints,
		DialTimeout: time.Second * 5,
	})
	if err != nil {
		panic(err)
	}
	// registrant actuator
	registrant := edclient.NewRegistrant(cli,
		config.Discovery.Etcd.Prefix,
		"/manager-center",
		// TODO get ip for network
		&edclient.NodeInfo{Server: "http://127.0.0.1:8082"},
		30,
	)

	err = registrant.Register(context.Background())
	if err != nil {
		panic(err)
	}

	err = server.Run(config.Serer.Port)
	if err != nil {
		panic(err)
	}
	registrant.Quit()
}
