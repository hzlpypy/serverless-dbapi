package server

import (
	"io"
	"net/http"
	"serverless-dbapi/pkg/exception"
	"serverless-dbapi/pkg/tool"
	"serverless-dbapi/pkg/valueobject"
	"sync"

	"github.com/gin-gonic/gin"
)

type Server interface {
	Run(addr ...string) error
}

var lock sync.Mutex
var ss *sharedServer

// shared server: when more than one services are in the same stand-alone
type sharedServer struct {
	server *gin.Engine
	lock   sync.Mutex
	isRun  bool
}

func (s *sharedServer) Run(addr ...string) error {
	s.lock.Lock()
	defer s.lock.Unlock()
	if !s.isRun {
		err := s.server.Run()
		if err == nil {
			s.isRun = true
		}
		return err
	}
	return nil
}

func newSharedServer() *sharedServer {
	lock.Lock()
	if ss == nil {
		ss = &sharedServer{
			server: gin.Default(),
			isRun:  false,
		}
	}
	defer lock.Unlock()
	return ss
}

// actuator server
func NewActuatorServer(function func(params *valueobject.Params) tool.Result[any]) Server {
	// impl by gin
	sharedServer := newSharedServer()
	sharedServer.server.POST("/actuator/api", func(ctx *gin.Context) {
		params, err := parseRequest(*ctx.Request, true)
		if err != nil {
			ctx.JSON(200, tool.ErrorResult[any](exception.PARSE_REQUEST_ERROR))
			return
		}
		result := function(params)
		ctx.JSON(200, tool.ResultToResponse(result))
	})
	return sharedServer.server
}

func NewManagerCenterServer(functions map[string]func(params *valueobject.Params) tool.Result[any]) Server {
	sharedServer := newSharedServer()
	sharedServer.server.Any("/manager-center/:resource", func(ctx *gin.Context) {
		params, err := parseRequest(*ctx.Request, true)
		if err != nil {
			ctx.JSON(200, tool.ErrorResult[any](exception.PARSE_REQUEST_ERROR))
			return
		}
		function := functions[tool.StringBuilder(ctx.Request.Method, "/", ctx.Param("resource"))]
		if function == nil {
			ctx.JSON(404, tool.ErrorResult[any](exception.HTTP_ROUTER_ERROR))
			return
		}
		result := function(params)
		ctx.JSON(200, tool.ResultToResponse(result))
	})
	return sharedServer.server
}

// parse request
func parseRequest(request http.Request, haveBody bool) (*valueobject.Params, error) {
	// get query params
	paramsQuery := request.URL.Query()

	params := &valueobject.Params{}
	params.QueryParams = paramsQuery

	// parse body
	if haveBody {
		body, err := io.ReadAll(request.Body)
		if err != nil {
			return nil, err
		}
		params.Body = body
	}

	return params, nil
}
