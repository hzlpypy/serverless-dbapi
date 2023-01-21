package store

import (
	"context"
	"encoding/json"
	"errors"
	"serverless-dbapi/pkg/cfg"
	"serverless-dbapi/pkg/entity"
	"serverless-dbapi/pkg/tool"
	"time"

	"github.com/google/uuid"
	"go.etcd.io/etcd/clientv3"
)

type Store interface {
	SaveDataBase(database entity.DatabaseConfig) error
	SaveApiGroup(apiGroup entity.ApiGroupConfig) error
	SaveApi(apiConfig entity.ApiConfig) error
	GetDataBases() ([]*entity.DatabaseConfig, error)
	GetApiGroups() ([]entity.ApiGroupConfig, error)
	GetApis(apiGroupId string) ([]entity.ApiConfig, error)
	GetApi(apiId string) (*entity.ApiConfig, error)
}

func StoreFactory(config cfg.StoreConfig) (Store, error) {
	if config.Etcd != nil {
		cli, err := clientv3.New(clientv3.Config{
			Endpoints:   config.Etcd.Endpoints,
			DialTimeout: time.Second * 5,
		})
		if err != nil {
			return nil, err
		}
		return &EtcdStore{
			client: cli,
			prefix: config.Etcd.Prefix,
		}, nil
	}
	return nil, errors.New("store only support etcd")
}

type EtcdStore struct {
	client *clientv3.Client
	prefix string
}

func (e *EtcdStore) SaveDataBase(database entity.DatabaseConfig) error {
	return e.commonPut(&database)
}

func (e *EtcdStore) SaveApiGroup(apiGroup entity.ApiGroupConfig) error {
	return e.commonPut(&apiGroup)
}

func (e *EtcdStore) SaveApi(apiConfig entity.ApiConfig) error {
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
		clientv3.OpPut(tool.StringBuilder(e.prefix, entity.API_PREFIX, apiConfig.ApiGroupId, "/", apiConfig.GetId()), string(dataByte)),
		clientv3.OpPut(tool.StringBuilder(e.prefix, entity.API_PREFIX, apiConfig.GetId()), string(dataByte)),
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

	_, err = e.client.Put(context.Background(), tool.StringBuilder(e.prefix, entity.GetPrefixId()), string(dataByte))
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

func (e *EtcdStore) GetDataBases() ([]*entity.DatabaseConfig, error) {
	resp, err := e.client.Get(context.Background(), tool.StringBuilder(e.prefix, entity.DATASOURCE_PREFIX), clientv3.WithPrefix())
	if err != nil {
		return nil, err
	}
	kv := resp.Kvs
	if len(kv) > 0 {
		result := make([]*entity.DatabaseConfig, 0, len(kv))
		for _, value := range kv {
			data := &entity.DatabaseConfig{}
			json.Unmarshal(value.Value, &data)
			result = append(result, data)
		}
		return result, nil
	}
	return []*entity.DatabaseConfig{}, nil
}

func (e *EtcdStore) GetApiGroups() ([]entity.ApiGroupConfig, error) {
	resp, err := e.client.Get(context.Background(), tool.StringBuilder(e.prefix, entity.API_GROUP_PREFIX), clientv3.WithPrefix())
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

func (e *EtcdStore) GetApis(apiGroupId string) ([]entity.ApiConfig, error) {
	resp, err := e.client.Get(context.Background(), tool.StringBuilder(e.prefix, entity.API_PREFIX, apiGroupId, "/"), clientv3.WithPrefix())
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

func (e *EtcdStore) GetApi(apiId string) (*entity.ApiConfig, error) {
	resp, err := e.client.Get(context.Background(), tool.StringBuilder(e.prefix, entity.API_PREFIX, apiId))
	if err != nil {
		return nil, err
	}
	kv := resp.Kvs
	if len(kv) > 0 {
		data := &entity.ApiConfig{}
		json.Unmarshal(kv[0].Value, &data)
		return data, nil
	}
	return nil, errors.New("api not found")
}
