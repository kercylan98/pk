// @APIVersion 1.0.0
// @Title School Schedule Plan Designer API
// @Description SSPD (School Schedule Plan Designer) Apis.
// @Contact kercylan@gmail.com
package routers

import (
	"pk/controllers/olds"
	"github.com/astaxie/beego"
)

func init() {
	// 主页面
	beego.Router("/", &olds.ViewController{}, "get:View")

	ns := beego.NewNamespace("/v1",
		beego.NSNamespace("/plan",
			beego.NSInclude(
				&olds.PlanController{},
			),
		),
		//beego.NSNamespace("/unit",
		//	beego.NSInclude(
		//		&controllers.UnitController{},
		//	),
		//),
		//beego.NSNamespace("/player",
		//	beego.NSInclude(
		//		&controllers.PlayerController{},
		//	),
		//),
	)

	beego.AddNamespace(ns)

}
