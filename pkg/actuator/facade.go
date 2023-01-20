package actuator

import (
	"bytes"
	"io"
	"net/http"
	"serverless-dbapi/pkg/tool"
	"serverless-dbapi/pkg/valueobject"

	"github.com/goccy/go-json"
)

// show the actuator server api
type ActuatorServer interface {
	ApiActuator(params *valueobject.Params) tool.Result[any]
}

// impl by http client
// TODO server LB
type HttpActuatorServer struct {
	client *http.Client
}

func NewHttpActuatorServer() ActuatorServer {
	return &HttpActuatorServer{
		client: &http.Client{},
	}
}

// TODO error-deal
func (h *HttpActuatorServer) ApiActuator(params *valueobject.Params) tool.Result[any] {
	url, _ := tool.BuildURL("http://localhost:8081", params.QueryParams)
	bodyBytes, _ := json.Marshal(params.Body)
	resp, _ := h.client.Post(url, "application/json", bytes.NewReader(bodyBytes))
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	data := &tool.Response[map[string]any]{}
	_ = json.Unmarshal(body, &data)
	return tool.SuccessResult(data)
}
