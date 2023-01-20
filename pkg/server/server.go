package server

import (
	"encoding/json"
	"io"
	"net/http"
	"serverless-dbapi/pkg/exception"
	"serverless-dbapi/pkg/mode"
	"serverless-dbapi/pkg/tool"
	"serverless-dbapi/pkg/valueobject"

	"github.com/gin-gonic/gin"
)

type Server interface {
	Run(addr ...string) error
}

type MemoryActuatorServer struct {
}

func (m *MemoryActuatorServer) Run(addr ...string) error {
	return nil
}

// actuator server
func NewActuatorServer(function func(params *valueobject.Params) tool.Result[any]) Server {
	if mode.MODE == mode.STANDALONE {
		return &MemoryActuatorServer{}
	} else {
		// impl by gin
		r := gin.Default()
		r.POST("/api", func(ctx *gin.Context) {
			params, err := parseRequest(*ctx.Request)
			if err != nil {
				ctx.JSON(exception.PARSE_REQUEST_ERROR.Code, exception.PARSE_REQUEST_ERROR.Msg)
			}
			result := function(params)
			ctx.JSON(200, tool.ResultToResponse(result))
		})
		return r
	}
}

// parse request
func parseRequest(request http.Request) (*valueobject.Params, error) {
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

	return &valueobject.Params{
		QueryParams: paramsQuery,
		Body:        paramsBody,
	}, nil
}
