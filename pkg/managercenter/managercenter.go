package managercenter

import (
	"serverless-dbapi/pkg/cfg"
	"serverless-dbapi/pkg/server"
	"serverless-dbapi/pkg/store"
	"serverless-dbapi/pkg/tool"
	"serverless-dbapi/pkg/valueobject"
	"strconv"
)

type ManagerCenter struct {
	store  store.Store
	server server.Server
}

func New(config *cfg.ManagerCenterConfig) (*ManagerCenter, error) {
	store, err := store.StoreFactory(*config.Store)
	if err != nil {
		return nil, err
	}
	handle := NewHandler(store)

	handleMap := make(map[string]func(params *valueobject.Params) tool.Result[any])
	handleMap["POST/database"] = handle.SaveDataBase
	handleMap["POST/api-group"] = handle.SaveApiGroup
	handleMap["POST/api"] = handle.SaveApi
	handleMap["GET/databases"] = handle.GetDataBases
	handleMap["GET/database"] = handle.GetDataBase
	handleMap["GET/api-groups"] = handle.GetApiGroups
	handleMap["GET/apis"] = handle.GetApis
	handleMap["GET/api"] = handle.GetApi

	server := server.NewManagerCenterServer(handleMap)
	return &ManagerCenter{
		store:  store,
		server: server,
	}, nil
}

func (m *ManagerCenter) Run(port int) error {
	m.server.Run(":" + strconv.Itoa(port))
	return nil
}
