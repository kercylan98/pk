package conn

const (
	CODE_SUCCESS = "code:0"								// 请求成功
	CODE_SYSTEM_ERR = "code:-1"							// 系统错误

	CODE_TOKEN_INVALID = "code:40000"					// 无效的access_token凭证
	CODE_TOKEN_UNAUTHORIZED = "code:40001"				// 尝试访问未授权的资源
	CODE_NOTFOUND = "code:40002"						// 尝试访问不存在的资源
	CODE_ILLEGAL_ARGS = "code:40003"					// 不合法的请求参数

	CODE_EXIST = "code:50000"							// 数据已存在

)

