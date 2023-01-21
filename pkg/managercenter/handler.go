package managercenter

import (
	"encoding/json"
	"serverless-dbapi/pkg/entity"
	"serverless-dbapi/pkg/exception"
	"serverless-dbapi/pkg/store"
	"serverless-dbapi/pkg/tool"
	"serverless-dbapi/pkg/valueobject"

	"github.com/go-playground/validator/v10"
)

type Handler struct {
	store store.Store
}

func NewHandler(store store.Store) Handler {
	return Handler{
		store: store,
	}
}

func (h *Handler) SaveDataBase(params *valueobject.Params) tool.Result[any] {
	database := &entity.DatabaseConfig{}
	err := json.Unmarshal(params.Body, &database)
	if err != nil {
		return tool.ErrorResult[any](exception.PARSE_REQUEST_ERROR)
	}
	valid := validator.New()
	if err := valid.Struct(database); err != nil {
		return tool.ErrorResult[any](exception.REQUIRE_PARAM, err.Error())
	}
	id, err := h.store.SaveDataBase(*database)
	if err != nil {
		return tool.SimpleErrorResult[any](500, err.Error())
	}
	return tool.SuccessResult[any](id)
}

func (h *Handler) SaveApiGroup(params *valueobject.Params) tool.Result[any] {
	apiGroup := &entity.ApiGroupConfig{}
	err := json.Unmarshal(params.Body, &apiGroup)
	if err != nil {
		return tool.ErrorResult[any](exception.PARSE_REQUEST_ERROR)
	}
	valid := validator.New()
	if err := valid.Struct(apiGroup); err != nil {
		return tool.ErrorResult[any](exception.REQUIRE_PARAM, err.Error())
	}
	id, err := h.store.SaveApiGroup(*apiGroup)
	if err != nil {
		return tool.SimpleErrorResult[any](500, err.Error())
	}
	return tool.SuccessResult[any](id)
}

func (h *Handler) SaveApi(params *valueobject.Params) tool.Result[any] {
	apiInfo := &entity.ApiConfig{}
	err := json.Unmarshal(params.Body, &apiInfo)
	if err != nil {
		return tool.ErrorResult[any](exception.PARSE_REQUEST_ERROR)
	}
	valid := validator.New()
	if err := valid.Struct(apiInfo); err != nil {
		return tool.ErrorResult[any](exception.REQUIRE_PARAM, err.Error())
	}
	id, err := h.store.SaveApi(*apiInfo)
	if err != nil {
		return tool.SimpleErrorResult[any](500, err.Error())
	}
	return tool.SuccessResult[any](id)
}

func (h *Handler) GetDataBases(params *valueobject.Params) tool.Result[any] {
	databases, err := h.store.GetDataBases()
	if err != nil {
		return tool.SimpleErrorResult[any](500, err.Error())
	}
	return tool.SuccessResult[any](databases)
}

func (h *Handler) GetApiGroups(params *valueobject.Params) tool.Result[any] {
	apiGroups, err := h.store.GetApiGroups()
	if err != nil {
		return tool.SimpleErrorResult[any](500, err.Error())
	}
	return tool.SuccessResult[any](apiGroups)
}

func (h *Handler) GetApis(params *valueobject.Params) tool.Result[any] {
	apis, err := h.store.GetApis(params.QueryParams["apiGroupId"][0])
	if err != nil {
		return tool.SimpleErrorResult[any](500, err.Error())
	}
	return tool.SuccessResult[any](apis)
}

func (h *Handler) GetApi(params *valueobject.Params) tool.Result[any] {
	api, err := h.store.GetApi(params.QueryParams["apiId"][0])
	if err != nil {
		return tool.SimpleErrorResult[any](500, err.Error())
	}
	return tool.SuccessResult[any](api)
}
