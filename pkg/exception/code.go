package exception

var PARSE_REQUEST_ERROR = New(10001, "解析请求参数错误")
var API_ID_IS_REQUIRE = New(10002, "apiId是必需的参数")
var REQUIRE_PARAM = New(10003, "缺少必要参数：%s")
var BUILD_URL_ERROR = New(10004, "生成url地址失败")
var RPC_ERROR = New(10005, "远程服务调用失败")
var RPC_RESPONSE_PARSE_ERROR = New(10006, "远程服务调用结果解析异常")
var DATASOURCE_NOT_FOUND = New(10007, "数据源路由错误")
