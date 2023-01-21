package managercenter

import (
	"errors"
	"serverless-dbapi/pkg/cfg"
	"serverless-dbapi/pkg/entity"
	"serverless-dbapi/pkg/mode"
)

// show manager center server api
type ManagerCenterServer interface {
	GetApiConfigByApiId(apiId string) (*entity.ApiConfig, error)
	GetDataBases() ([]*entity.DatabaseConfig, error)
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
func (m *MockManagerCenterServer) GetApiConfigByApiId(apiId string) (*entity.ApiConfig, error) {
	return &entity.ApiConfig{
		Id:           apiId,
		ApiGroupId:   "1",
		ApiType:      entity.DATA_API_TYPE,
		Sql:          "select * from tb_tmp01 where id = ? and deptId = ?",
		ParamKey:     []string{"id", "deptId"},
		DataSourceId: "1",
	}, nil
}

func (m *MockManagerCenterServer) GetDataBases() ([]*entity.DatabaseConfig, error) {
	return []*entity.DatabaseConfig{
		{
			Id:         "1",
			Name:       "test",
			DriverName: "mysql",
			Url:        "root:123456@tcp(localhost:3306)/test?charset=utf8",
		},
	}, nil
}

type MemoryManagerCenterServer struct {
	managerCenter *ManagerCenter
}

func (m *MemoryManagerCenterServer) GetApiConfigByApiId(apiId string) (*entity.ApiConfig, error) {
	return m.managerCenter.store.GetApi(apiId)
}

func (m *MemoryManagerCenterServer) GetDataBases() ([]*entity.DatabaseConfig, error) {
	return m.managerCenter.store.GetDataBases()
}

type HttpManagerCenterServer struct {
}

func (h *HttpManagerCenterServer) GetApiConfigByApiId(apiId string) (*entity.ApiConfig, error) {
	return nil, nil
}

func (h *HttpManagerCenterServer) GetDataBases() ([]*entity.DatabaseConfig, error) {
	return nil, nil
}

func NewManagerCenterServer(config *cfg.ManagerCenterConfig) (ManagerCenterServer, error) {
	if mode.MODE == mode.MOCK {
		return &MockManagerCenterServer{}, nil
	}
	if mode.MODE == mode.STANDALONE {
		managerCenter, err := New(config)
		if err != nil {
			return nil, err
		}
		return &MemoryManagerCenterServer{
			managerCenter: managerCenter,
		}, nil
	}
	if mode.MODE == mode.CLUSTER {
		return &HttpManagerCenterServer{}, nil
	}
	return nil, errors.New("manager center server not found")
}
