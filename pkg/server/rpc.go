package server

import "serverless-db/pkg/entity"

type ManagerCenterServer interface {
	GetApiConfigByApiId(apiId string) entity.ApiConfig
}

type MockManagerCenterServer struct {
}

func (m *MockManagerCenterServer) GetApiConfigByApiId(apiId string) entity.ApiConfig {
	return entity.ApiConfig{
		ApiId:    apiId,
		ApiType:  entity.DATA_API_TYPE,
		Sql:      "select * from tb_tmp01 where id = ? and deptId = ?",
		ParamKey: []string{"id", "deptId"},
	}
}
