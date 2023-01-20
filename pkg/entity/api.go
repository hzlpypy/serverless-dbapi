package entity

// api type
const (
	DATA_API_TYPE = 1
)

// api config info
type ApiConfig struct {
	ApiId    string
	ApiType  int
	Sql      string
	ParamKey []string
}
