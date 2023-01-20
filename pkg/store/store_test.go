package store

import (
	"fmt"
	"serverless-dbapi/pkg/entity"
	"testing"
	"time"

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
	etcdStore.saveDataBase(entity.DatabaseConfig{
		Id:         "1",
		Name:       "datasource-1",
		DriverName: "mysql",
		Url:        "root:123456@tcp(localhost:3306)/test?charset=utf8",
	})
	etcdStore.saveApiGroup(entity.ApiGroupConfig{
		Id:   "1",
		Name: "group-1",
	})
	etcdStore.saveApiGroup(entity.ApiGroupConfig{
		Id:   "2",
		Name: "group-2",
	})
	etcdStore.saveApi(entity.ApiConfig{
		Id:           "1",
		ApiGroupId:   "1",
		ApiType:      entity.DATA_API_TYPE,
		Sql:          "select * from tb_tmp01 where id = ? and deptId = ?",
		ParamKey:     []string{"id", "deptId"},
		DataSourceId: "1",
	})
	etcdStore.saveApi(entity.ApiConfig{
		Id:           "2",
		ApiGroupId:   "1",
		ApiType:      entity.DATA_API_TYPE,
		Sql:          "select * from tb_tmp01 where id = ? and deptId = ?",
		ParamKey:     []string{"id", "deptId"},
		DataSourceId: "1",
	})

	datasources, err := etcdStore.getDataBases()
	if err != nil {
		t.Error(err)
	}
	fmt.Println(datasources)

	apiGroups, err := etcdStore.getApiGroups()
	if err != nil {
		t.Error(err)
	}
	fmt.Println(apiGroups)

	apis, err := etcdStore.getApis("1")
	if err != nil {
		t.Error(err)
	}
	fmt.Println(apis)

	api, err := etcdStore.getApi("1")
	if err != nil {
		t.Error(err)
	}
	fmt.Println(api)
}
