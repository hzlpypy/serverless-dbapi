package store

import (
	"serverless-dbapi/pkg/entity"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.etcd.io/etcd/clientv3"
)

func Test_Store(t *testing.T) {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"127.0.0.1:2379"},
		DialTimeout: time.Second * 5,
	})
	if err != nil {
		t.Error(err)
	}

	etcdStore := &EtcdStore{}
	etcdStore.client = cli
	etcdStore.prefix = "/test-store"
	database := entity.DatabaseConfig{
		Id:         "1",
		Name:       "datasource-1",
		DriverName: "mysql",
		Url:        "root:123456@tcp(localhost:3306)/test?charset=utf8",
	}
	etcdStore.SaveDataBase(database)
	apiGroupOne := entity.ApiGroupConfig{
		Id:   "1",
		Name: "group-1",
	}
	apiGroupTwo := entity.ApiGroupConfig{
		Id:   "2",
		Name: "group-2",
	}
	etcdStore.SaveApiGroup(apiGroupOne)
	etcdStore.SaveApiGroup(apiGroupTwo)

	apiOne := entity.ApiConfig{
		Id:           "1",
		ApiGroupId:   "1",
		ApiType:      entity.DATA_API_TYPE,
		Sql:          "select * from tb_tmp01 where id = ? and deptId = ?",
		ParamKey:     []string{"id", "deptId"},
		DataSourceId: "1",
	}
	apiTwo := entity.ApiConfig{
		Id:           "2",
		ApiGroupId:   "1",
		ApiType:      entity.DATA_API_TYPE,
		Sql:          "select * from tb_tmp01 where id = ? and deptIds in (?)",
		ParamKey:     []string{"id", "deptIds"},
		DataSourceId: "1",
	}
	etcdStore.SaveApi(apiOne)
	etcdStore.SaveApi(apiTwo)

	datasources, err := etcdStore.GetDataBases()
	if err != nil {
		t.Error(err)
	}
	if len(datasources) != 1 {
		t.Error("datasources count error")
	}
	for _, value := range datasources {
		assert.Equal(t, value, database)
	}

	apiGroups, err := etcdStore.GetApiGroups()
	if err != nil {
		t.Error(err)
	}
	if len(apiGroups) != 2 {
		t.Error("api group count error")
	}
	for _, value := range apiGroups {
		if value.Id == "1" {
			assert.Equal(t, value, apiGroupOne)
		} else if value.Id == "2" {
			assert.Equal(t, value, apiGroupTwo)
		} else {
			t.Error("api group not exist")
		}
	}

	apis, err := etcdStore.GetApis("1")
	if err != nil {
		t.Error(err)
	}
	if len(apis) != 2 {
		t.Error("api count error")
	}
	for _, value := range apis {
		if value.Id == "1" {
			assert.Equal(t, value, apiOne)
		} else if value.Id == "2" {
			assert.Equal(t, value, apiTwo)
		} else {
			t.Error("api not exist")
		}
	}

	api, err := etcdStore.GetApi("1")
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, *api, apiOne)
}
