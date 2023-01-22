package store

import (
	"context"
	"encoding/json"
	"errors"
	"serverless-dbapi/pkg/cfg"
	"serverless-dbapi/pkg/entity"
	"serverless-dbapi/pkg/tool"
	"serverless-dbapi/pkg/valueobject"
	"time"

	"github.com/google/uuid"
	"go.etcd.io/etcd/clientv3"
)

type Store interface {
	SaveDataBase(database entity.DatabaseConfig) (string, error)
	SaveApiGroup(apiGroup entity.ApiGroupConfig) (string, error)
	SaveApi(apiConfig entity.ApiConfig) (string, error)
	GetAllDataBases() ([]*entity.DatabaseConfig, error)
	GetDataBases(page valueobject.Cursor) ([]*entity.DatabaseConfig, error)
	GetApiGroups(page valueobject.Cursor) ([]entity.ApiGroupConfig, error)
	GetApis(apiGroupId string, page valueobject.Cursor) ([]entity.ApiConfig, error)
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
		prefix := config.Etcd.Prefix
		return &EtcdStore{
			client:             cli,
			prefix:             prefix,
			saveApiPrefix:      tool.StringBuilder(prefix, entity.API_PREFIX),
			saveDatabasePrefix: tool.StringBuilder(prefix, entity.DATASOURCE_PREFIX),
			saveApiGroupPrefix: tool.StringBuilder(prefix, entity.API_GROUP_PREFIX),
		}, nil
	}
	return nil, errors.New("store only support etcd")
}

type EtcdStore struct {
	client *clientv3.Client
	prefix string

	saveApiPrefix      string
	saveDatabasePrefix string
	saveApiGroupPrefix string
}

func (e *EtcdStore) SaveDataBase(database entity.DatabaseConfig) (string, error) {
	return e.commonPut(&database, e.saveDatabasePrefix)
}

func (e *EtcdStore) SaveApiGroup(apiGroup entity.ApiGroupConfig) (string, error) {
	return e.commonPut(&apiGroup, e.saveApiGroupPrefix)
}

func (e *EtcdStore) SaveApi(apiConfig entity.ApiConfig) (string, error) {
	err := fillId(&apiConfig)
	if err != nil {
		return "", err
	}
	dataByte, err := json.Marshal(apiConfig)
	if err != nil {
		return "", err
	}

	txn := e.client.Txn(context.Background())
	_, err = txn.If().Then(
		clientv3.OpPut(tool.StringBuilder(e.saveApiPrefix, apiConfig.ApiGroupId, "/", apiConfig.GetId()), string(dataByte)),
		clientv3.OpPut(tool.StringBuilder(e.saveApiPrefix, apiConfig.GetId()), string(dataByte)),
	).Else().Commit()
	if err != nil {
		return "", err
	}
	return apiConfig.Id, nil
}

func (e *EtcdStore) commonPut(entity entity.IdCommon, prefix string) (string, error) {
	err := fillId(entity)
	if err != nil {
		return "", err
	}
	dataByte, err := json.Marshal(entity)
	if err != nil {
		return "", err
	}

	_, err = e.client.Put(context.Background(), tool.StringBuilder(prefix, entity.GetId()), string(dataByte))
	if err != nil {
		return "", err
	}

	return entity.GetId(), nil
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

func (e *EtcdStore) GetAllDataBases() ([]*entity.DatabaseConfig, error) {
	resp, err := e.client.Get(
		context.Background(),
		tool.StringBuilder(e.saveDatabasePrefix),
		clientv3.WithPrefix(),
	)
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

func (e *EtcdStore) GetDataBases(page valueobject.Cursor) ([]*entity.DatabaseConfig, error) {
	specialPage(&page)
	resp, err := e.client.Get(
		context.Background(),
		tool.StringBuilder(e.saveDatabasePrefix, page.Continue),
		clientv3.WithRange(tool.StringBuilder(e.prefix, "/datasource0")),
		clientv3.WithSort(clientv3.SortByKey, clientv3.SortAscend),
		clientv3.WithLimit(int64(page.Limit)),
	)
	if err != nil {
		return nil, err
	}
	kv := resp.Kvs

	if len(kv) > 0 {
		result := make([]*entity.DatabaseConfig, 0, len(kv))
		for _, value := range kv {
			data := &entity.DatabaseConfig{}
			json.Unmarshal(value.Value, &data)
			if data.Id != page.Continue {
				result = append(result, data)
			}
		}
		return result, nil
	}
	return []*entity.DatabaseConfig{}, nil
}

func (e *EtcdStore) GetApiGroups(page valueobject.Cursor) ([]entity.ApiGroupConfig, error) {
	specialPage(&page)
	resp, err := e.client.Get(
		context.Background(),
		tool.StringBuilder(e.saveApiGroupPrefix, page.Continue),
		clientv3.WithRange(tool.StringBuilder(e.prefix, "/api-group0")),
		clientv3.WithSort(clientv3.SortByKey, clientv3.SortAscend),
		clientv3.WithLimit(int64(page.Limit)),
	)
	if err != nil {
		return nil, err
	}
	kv := resp.Kvs
	if len(kv) > 0 {
		result := make([]entity.ApiGroupConfig, 0, len(kv))
		for _, value := range kv {
			data := &entity.ApiGroupConfig{}
			json.Unmarshal(value.Value, &data)
			if data.Id != page.Continue {
				result = append(result, *data)
			}
		}
		return result, nil
	}
	return []entity.ApiGroupConfig{}, nil
}

func (e *EtcdStore) GetApis(apiGroupId string, page valueobject.Cursor) ([]entity.ApiConfig, error) {
	specialPage(&page)
	resp, err := e.client.Get(
		context.Background(),
		tool.StringBuilder(e.saveApiPrefix, apiGroupId, "/", page.Continue),
		clientv3.WithRange(tool.StringBuilder(e.prefix, "/api0")),
		clientv3.WithSort(clientv3.SortByKey, clientv3.SortAscend),
		clientv3.WithLimit(int64(page.Limit)),
	)
	if err != nil {
		return nil, err
	}
	kv := resp.Kvs
	if len(kv) > 0 {
		result := make([]entity.ApiConfig, 0, len(kv))
		for _, value := range kv {
			data := &entity.ApiConfig{}
			json.Unmarshal(value.Value, &data)
			if data.Id != page.Continue {
				result = append(result, *data)
			}
		}
		return result, nil
	}
	return []entity.ApiConfig{}, nil
}

func (e *EtcdStore) GetApi(apiId string) (*entity.ApiConfig, error) {
	resp, err := e.client.Get(
		context.Background(),
		tool.StringBuilder(e.saveApiPrefix, apiId),
	)
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

func specialPage(pageInfo *valueobject.Cursor) {
	if pageInfo.Continue != "" {
		pageInfo.Limit += 1
	}
}
