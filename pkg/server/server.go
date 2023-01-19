package server

import (
	"net/http"
	"serverless-db/pkg/result"

	"github.com/gin-gonic/gin"
)

// server
// impl by gin
type Server interface {
	Run(addr ...string) error
}

// actuator server
func NewActuatorServer(function func(request http.Request) result.Result[any]) Server {
	r := gin.Default()
	r.POST("/api", func(ctx *gin.Context) {
		result := function(*ctx.Request)
		if result.IsError() {
			err := result.Err
			ctx.JSON(err.Code, err.Msg)
		} else {
			ctx.JSON(200, result.Data)
		}
	})
	return r
}
