package main

import (
	"context"
	"errors"
	"os"
	"serverless-dbapi/cmd/mode"
	"serverless-dbapi/pkg/actuator"
	"serverless-dbapi/pkg/cfg"
	"serverless-dbapi/pkg/managercenter"
	"time"

	"github.com/coreos/etcd/clientv3"
	_ "github.com/go-sql-driver/mysql"
	edclient "github.com/kiraqjx/ed-client"
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
	managerCenterServer, err := newManagerCenterServer(config)
	if err != nil {
		panic(err)
	}

	if mode.MODE != mode.CLUSTER {
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

func newManagerCenterServer(config cfg.Config) (managercenter.ManagerCenterServer, error) {
	if mode.MODE == mode.MOCK {
		return managercenter.NewMockManagerCenterServer()
	}
	if mode.MODE == mode.STANDALONE {
		return managercenter.NewMemoryManagerCenterServer(config.ManagerCenter)
	}
	if mode.MODE == mode.CLUSTER {
		cli, err := clientv3.New(clientv3.Config{
			Endpoints:   config.Discovery.Etcd.Endpoints,
			DialTimeout: time.Second * 5,
		})
		if err != nil {
			return nil, err
		}
		// registrant actuator
		registrant := edclient.NewRegistrant(cli,
			config.Discovery.Etcd.Prefix,
			"/actuator",
			// TODO get ip for network
			&edclient.NodeInfo{Server: "http://127.0.0.1:8081"},
			30,
		)
		err = registrant.Register(context.Background())
		if err != nil {
			return nil, err
		}

		// watch manager center
		watcher := edclient.NewWatcher(cli, config.Discovery.Etcd.Prefix, "/manager-center")
		err = watcher.Start(context.Background())
		if err != nil {
			registrant.Quit()
			return nil, err
		}

		http, err := managercenter.NewHttpManagerCenterServer(edclient.NewLbFromMap(watcher.Nodes))
		if err != nil {
			registrant.Quit()
			return nil, err
		}

		go func() {
			for {
				<-watcher.ChangeEvent()
				http.Lb.ChangeNodesFromMap(watcher.Nodes)
			}
		}()

		return http, nil
	}
	return nil, errors.New("manager center server not found")
}
