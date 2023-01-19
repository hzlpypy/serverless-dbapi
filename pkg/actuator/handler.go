package actuator

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"serverless-db/pkg/entity"
	"serverless-db/pkg/exception"
	"serverless-db/pkg/result"
	"serverless-db/pkg/server"
)

const API_ID_FROM_QUERY = "apiId"

type Handler struct {
	dbConn        *sql.DB
	managerCenter server.ManagerCenterServer
}

func NewHandle(dbConn *sql.DB, managerCenter server.ManagerCenterServer) Handler {
	return Handler{
		dbConn:        dbConn,
		managerCenter: managerCenter,
	}
}

// common handler
func (h *Handler) Handler(request http.Request) result.Result[any] {
	params, err := parseRequest(request)
	if err != nil {
		return result.ErrorResult(exception.PARSE_REQUEST_ERROR)
	}
	apiIds := params.queryParams[API_ID_FROM_QUERY]
	if len(apiIds) != 1 {
		return result.ErrorResult(exception.API_ID_IS_REQUIRE)
	}
	apiId := apiIds[0]
	apiConfig := h.managerCenter.GetApiConfigByApiId(apiId)
	return h.exec(apiConfig, *params)
}

type params struct {
	headers     map[string][]string
	queryParams map[string][]string
	body        map[string]any
}

// parse request
func parseRequest(request http.Request) (*params, error) {
	// get header params
	paramsHeader := request.Header

	// get query params
	paramsQuery := request.URL.Query()

	// parse body
	paramsBody := make(map[string]any)
	s, _ := io.ReadAll(request.Body)
	if len(s) > 0 {
		err := json.Unmarshal(s, &paramsBody)
		if err != nil {
			return nil, err
		}
	}

	return &params{
		headers:     paramsHeader,
		queryParams: paramsQuery,
		body:        paramsBody,
	}, nil
}

func (h *Handler) exec(apiConfig entity.ApiConfig, params params) result.Result[any] {
	args := make([]any, len(apiConfig.ParamKey))
	for index, value := range apiConfig.ParamKey {
		if value, ok := params.body[value]; ok {
			args[index] = value
		} else {
			return result.ErrorResult(exception.REQUIRE_PARAM)
		}
	}
	rows, err := h.dbConn.Query(apiConfig.Sql, args...)
	if err != nil {
		fmt.Println(err)
		return result.SimpleErrorResult(500, err.Error())
	}
	defer rows.Close()

	columns, _ := rows.Columns()
	columnLength := len(columns)
	cache := make([]any, columnLength)
	for index := range cache {
		var a any
		cache[index] = &a
	}
	var list []map[string]any
	for rows.Next() {
		_ = rows.Scan(cache...)
		item := make(map[string]any)
		for i, data := range cache {
			item[columns[i]] = data
		}
		list = append(list, item)
	}

	return result.SuccessResult(list)
}
