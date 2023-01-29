package managercenter

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"serverless-dbapi/pkg/cfg"
	"serverless-dbapi/pkg/entity"
	"serverless-dbapi/pkg/tool"
	"serverless-dbapi/pkg/valueobject"

	edclient "github.com/kiraqjx/ed-client"
)

// show manager center server api
type ManagerCenterServer interface {
	GetApiConfigByApiId(apiId string) (*entity.ApiConfig, error)
	GetDataBases() ([]*entity.DatabaseConfig, error)
	GetDataBase(dataBaseId string) (*entity.DatabaseConfig, error)
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

func (m *MockManagerCenterServer) GetDataBase(_ string) (*entity.DatabaseConfig, error) {
	return &entity.DatabaseConfig{
		Id:         "1",
		Name:       "test",
		DriverName: "mysql",
		Url:        "root:123456@tcp(localhost:3306)/test?charset=utf8",
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

func (m *MemoryManagerCenterServer) GetDataBase(dataBaseId string) (*entity.DatabaseConfig, error) {
	return m.managerCenter.store.GetDataBase(dataBaseId)
}

type HttpManagerCenterServer struct {
	Lb     *edclient.Lb
	client *http.Client
}

func NewHttpManagerCenterServer(lb *edclient.Lb) (*HttpManagerCenterServer, error) {
	return &HttpManagerCenterServer{
		Lb: lb,
	}, nil
}

func (h *HttpManagerCenterServer) GetApiConfigByApiId(apiId string) (*entity.ApiConfig, error) {
	node := h.Lb.Lb()
	params := &valueobject.Params{}
	params.QueryParams["apiId"] = []string{apiId}
	url, err := tool.BuildURL(tool.StringBuilder(node.Server, "/manager-center/api"), params.QueryParams)
	if err != nil {
		return nil, err
	}
	resp, err := h.client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	data := &tool.Response[*entity.ApiConfig]{}
	err = json.Unmarshal(body, &data)
	if err != nil {
		return nil, err
	}
	return data.Data, nil
}

func (h *HttpManagerCenterServer) GetDataBases() ([]*entity.DatabaseConfig, error) {
	return nil, errors.New("not support")
}

func (h *HttpManagerCenterServer) GetDataBase(dataBaseId string) (*entity.DatabaseConfig, error) {
	node := h.Lb.Lb()
	params := &valueobject.Params{}
	params.QueryParams["id"] = []string{dataBaseId}
	url, err := tool.BuildURL(tool.StringBuilder(node.Server, "/manager-center/database"), params.QueryParams)
	if err != nil {
		return nil, err
	}
	resp, err := h.client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	data := &tool.Response[*entity.DatabaseConfig]{}
	err = json.Unmarshal(body, &data)
	if err != nil {
		return nil, err
	}
	return data.Data, nil
}
