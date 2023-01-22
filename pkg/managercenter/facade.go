package managercenter

import (
	"serverless-dbapi/pkg/cfg"
	"serverless-dbapi/pkg/entity"
)

// show manager center server api
type ManagerCenterServer interface {
	GetApiConfigByApiId(apiId string) (*entity.ApiConfig, error)
	GetDataBases() ([]*entity.DatabaseConfig, error)
}

// manager center server impl mock
type MockManagerCenterServer struct {
}

func NewMockManagerCenterServer() (ManagerCenterServer, error) {
	return &MockManagerCenterServer{}, nil
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

func NewMemoryManagerCenterServer(config *cfg.ManagerCenterConfig) (ManagerCenterServer, error) {
	managerCenter, err := New(config)
	if err != nil {
		return nil, err
	}
	return &MemoryManagerCenterServer{
		managerCenter: managerCenter,
	}, nil
}

func (m *MemoryManagerCenterServer) GetApiConfigByApiId(apiId string) (*entity.ApiConfig, error) {
	return m.managerCenter.store.GetApi(apiId)
}

func (m *MemoryManagerCenterServer) GetDataBases() ([]*entity.DatabaseConfig, error) {
	return m.managerCenter.store.GetAllDataBases()
}

type HttpManagerCenterServer struct {
}

func NewHttpManagerCenterServer() (ManagerCenterServer, error) {
	return &HttpManagerCenterServer{}, nil
}

func (h *HttpManagerCenterServer) GetApiConfigByApiId(apiId string) (*entity.ApiConfig, error) {
	return nil, nil
}

func (h *HttpManagerCenterServer) GetDataBases() ([]*entity.DatabaseConfig, error) {
	return nil, nil
}
