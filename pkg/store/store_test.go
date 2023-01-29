package store

import (
	"context"
	"serverless-dbapi/pkg/entity"
	"serverless-dbapi/pkg/tool"
	"serverless-dbapi/pkg/valueobject"
	"testing"
	"time"

	"github.com/coreos/etcd/clientv3"
	"github.com/stretchr/testify/assert"
)

func Test_Store(t *testing.T) {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"127.0.0.1:2379"},
		DialTimeout: time.Second * 5,
	})
	if err != nil {
		t.Error(err)
	}

	// clean data
	cli.Delete(context.Background(), "/test-store", clientv3.WithPrefix())

	defer func() {
		// clean data
		cli.Delete(context.Background(), "/test-store", clientv3.WithPrefix())
	}()

	prefix := "/test-store"
	etcdStore := &EtcdStore{
		client:             cli,
		prefix:             prefix,
		saveApiPrefix:      tool.StringBuilder(prefix, entity.API_PREFIX),
		saveDatabasePrefix: tool.StringBuilder(prefix, entity.DATASOURCE_PREFIX),
		saveApiGroupPrefix: tool.StringBuilder(prefix, entity.API_GROUP_PREFIX),
	}
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

	// all database
	datasources, err := etcdStore.GetAllDataBases()
	if err != nil {
		t.Error(err)
	}
	if len(datasources) != 1 {
		t.Error("datasources count error")
	}
	for _, value := range datasources {
		assert.Equal(t, *value, database)
	}

	// all database
	datasources, err = etcdStore.GetDataBases(valueobject.Cursor{Continue: "", Limit: 1})
	if err != nil {
		t.Error(err)
	}
	if len(datasources) != 1 {
		t.Error("datasources count error")
	}
	for _, value := range datasources {
		assert.Equal(t, *value, database)
	}

	// single database
	datasource, err := etcdStore.GetDataBase("1")
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, *datasource, database)

	// all api groups
	apiGroups, err := etcdStore.GetApiGroups(valueobject.Cursor{Continue: "", Limit: 2})
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

	// api group page
	apiGroups, err = etcdStore.GetApiGroups(valueobject.Cursor{Continue: "1", Limit: 3})
	if err != nil {
		t.Error(err)
	}
	if len(apiGroups) != 1 {
		t.Error("api group count error")
	}
	for _, value := range apiGroups {
		if value.Id == "2" {
			assert.Equal(t, value, apiGroupTwo)
		} else {
			t.Error("api group not exist")
		}
	}

	// all api by api group id
	apis, err := etcdStore.GetApis("1", valueobject.Cursor{Continue: "", Limit: 2})
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

	// api by api id
	api, err := etcdStore.GetApi("1")
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, *api, apiOne)
}
