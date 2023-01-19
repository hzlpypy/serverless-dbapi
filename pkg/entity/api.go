package entity

const (
	DATA_API_TYPE = 1
)

type ApiConfig struct {
	ApiId    string
	ApiType  int
	Sql      string
	ParamKey []string
}
