package exception

var PARSE_REQUEST_ERROR = New(10001, "解析请求参数错误")
var API_ID_IS_REQUIRE = New(10002, "apiId是必需的参数")
var REQUIRE_PARAM = New(10003, "缺少必要参数")
