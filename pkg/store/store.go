package store

import (
	"context"
	"encoding/json"
	"serverless-dbapi/pkg/entity"

	"github.com/google/uuid"
	"go.etcd.io/etcd/clientv3"
)

type Store interface {
	saveDataBase(database entity.DatabaseConfig) error
	saveApiGroup(apiGroup entity.ApiGroupConfig) error
	saveApi(apiConfig entity.ApiConfig) error
	getDataBases() ([]entity.DatabaseConfig, error)
	getApiGroups() ([]entity.ApiGroupConfig, error)
	getApis(apiGroupId string) ([]entity.ApiConfig, error)
	getApi(apiId string) ([]entity.ApiConfig, error)
}

type EtcdStore struct {
	client *clientv3.Client
}

func (e *EtcdStore) saveDataBase(database entity.DatabaseConfig) error {
	return e.commonPut(&database)
}

func (e *EtcdStore) saveApiGroup(apiGroup entity.ApiGroupConfig) error {
	return e.commonPut(&apiGroup)
}

func (e *EtcdStore) saveApi(apiConfig entity.ApiConfig) error {
	err := fillId(&apiConfig)
	if err != nil {
		return err
	}
	dataByte, err := json.Marshal(apiConfig)
	if err != nil {
		return err
	}

	txn := e.client.Txn(context.Background())
	_, err = txn.If().Then(
		clientv3.OpPut(entity.API_PREFIX+apiConfig.ApiGroupId+"/"+apiConfig.GetId(), string(dataByte)),
		clientv3.OpPut(entity.API_PREFIX+apiConfig.GetId(), string(dataByte)),
	).Else().Commit()
	if err != nil {
		return err
	}
	return nil
}

func (e *EtcdStore) commonPut(entity entity.IdCommon) error {
	err := fillId(entity)
	if err != nil {
		return err
	}
	dataByte, err := json.Marshal(entity)
	if err != nil {
		return err
	}

	_, err = e.client.Put(context.Background(), entity.GetPrefixId(), string(dataByte))
	if err != nil {
		return err
	}

	return nil
}

func fillId(entity entity.IdCommon) error {
	if entity.GetId() == "" {
		uuid, err := uuid.NewUUID()
		if err != nil {
			return err
		}
		entity.SetId(uuid.String())
	}
	return nil
}

func (e *EtcdStore) getDataBases() ([]entity.DatabaseConfig, error) {
	resp, err := e.client.Get(context.Background(), entity.DATASOURCE_PREFIX, clientv3.WithPrefix())
	if err != nil {
		return nil, err
	}
	kv := resp.Kvs
	if len(kv) > 0 {
		result := make([]entity.DatabaseConfig, 0, len(kv))
		for _, value := range kv {
			data := &entity.DatabaseConfig{}
			json.Unmarshal(value.Value, &data)
			result = append(result, *data)
		}
		return result, nil
	}
	return []entity.DatabaseConfig{}, nil
}

func (e *EtcdStore) getApiGroups() ([]entity.ApiGroupConfig, error) {
	resp, err := e.client.Get(context.Background(), entity.API_GROUP_PREFIX, clientv3.WithPrefix())
	if err != nil {
		return nil, err
	}
	kv := resp.Kvs
	if len(kv) > 0 {
		result := make([]entity.ApiGroupConfig, 0, len(kv))
		for _, value := range kv {
			data := &entity.ApiGroupConfig{}
			json.Unmarshal(value.Value, &data)
			result = append(result, *data)
		}
		return result, nil
	}
	return []entity.ApiGroupConfig{}, nil
}

func (e *EtcdStore) getApis(apiGroupId string) ([]entity.ApiConfig, error) {
	resp, err := e.client.Get(context.Background(), entity.API_PREFIX+apiGroupId+"/", clientv3.WithPrefix())
	if err != nil {
		return nil, err
	}
	kv := resp.Kvs
	if len(kv) > 0 {
		result := make([]entity.ApiConfig, 0, len(kv))
		for _, value := range kv {
			data := &entity.ApiConfig{}
			json.Unmarshal(value.Value, &data)
			result = append(result, *data)
		}
		return result, nil
	}
	return []entity.ApiConfig{}, nil
}

func (e *EtcdStore) getApi(apiId string) ([]entity.ApiConfig, error) {
	resp, err := e.client.Get(context.Background(), entity.API_PREFIX+apiId)
	if err != nil {
		return nil, err
	}
	kv := resp.Kvs
	if len(kv) > 0 {
		result := make([]entity.ApiConfig, 0, len(kv))
		for _, value := range kv {
			data := &entity.ApiConfig{}
			json.Unmarshal(value.Value, &data)
			result = append(result, *data)
		}
		return result, nil
	}
	return []entity.ApiConfig{}, nil
}
