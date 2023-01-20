package managercenter

import (
	"serverless-dbapi/pkg/entity"
	"serverless-dbapi/pkg/mode"
)

// show manager center server api
type ManagerCenterServer interface {
	GetApiConfigByApiId(apiId string) entity.ApiConfig
}

// manager center server impl mock
type MockManagerCenterServer struct {
}

// mysql
// create table tb_tmp01
// (
// id INT(11),
// name VARCHAR(25),
// deptId INT(11),
// salary FLOAT
// );
func (m *MockManagerCenterServer) GetApiConfigByApiId(apiId string) entity.ApiConfig {
	return entity.ApiConfig{
		Id:           apiId,
		ApiGroupId:   "1",
		ApiType:      entity.DATA_API_TYPE,
		Sql:          "select * from tb_tmp01 where id = ? and deptId = ?",
		ParamKey:     []string{"id", "deptId"},
		DataSourceId: "1",
	}
}

func NewManagerCenterServer() ManagerCenterServer {
	if mode.MODE == mode.MOCK {
		return &MockManagerCenterServer{}
	}
	return nil
}
