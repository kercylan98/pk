package basic

import (
	"pk/app"
	"github.com/astaxie/beego"
)

const (
	private_KEY_ENGINES = "ENGINES"	// 获取所有引擎的键
)

// 所有控制器的基类
type Controller struct {
	beego.Controller
}


// 执行具体函数前初始化
func (slf *Controller) Prepare() {
	// 初始化赋值
	slf.Data[private_KEY_ENGINES] = app.Get()
	// 身份校验
	//result := conn.Dispose(func() (interface{}, interface{}) {
	//	token, err := request.ParseFromRequest(slf.Ctx.Request, request.AuthorizationHeaderExtractor,
	//		func(token *jwt.Token) (interface{}, error) {
	//			return []byte(engines.SECRET_KEY), nil
	//		})
	//
	//	if err == nil {
	//		if token.Valid {
	//			// 通过
	//		} else {
	//			slf.Ctx.Output.SetStatus(http.StatusUnauthorized)
	//			return conn.CODE_TOKEN_INVALID, "err:无效的身份令牌凭证"
	//		}
	//	} else {
	//		logs.Error(err)
	//		slf.Ctx.Output.SetStatus(http.StatusUnauthorized)
	//		return conn.CODE_TOKEN_UNAUTHORIZED, "err:无法对未授权的资源进行访问"
	//	}
	//	return nil, nil
	//})
	//if result.Code != 200 {
	//	slf.Data["json"] = result
	//	slf.ServeJSON()
	//	slf.StopRun()
	//}
}

// 获取应用程序引擎
func (slf *Controller) App() *app.Application {
	return slf.Data[private_KEY_ENGINES].(*app.Application)
}