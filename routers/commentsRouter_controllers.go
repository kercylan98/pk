package routers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context/param"
)

func init() {

    beego.GlobalControllerRouter["pk/ibps/controllers:PlayerController"] = append(beego.GlobalControllerRouter["pk/ibps/controllers:PlayerController"],
        beego.ControllerComments{
            Method: "Post",
            Router: `/`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["pk/ibps/controllers:UnitController"] = append(beego.GlobalControllerRouter["pk/ibps/controllers:UnitController"],
        beego.ControllerComments{
            Method: "Post",
            Router: `/`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["pk/ibps/controllers:UnitController"] = append(beego.GlobalControllerRouter["pk/ibps/controllers:UnitController"],
        beego.ControllerComments{
            Method: "Get",
            Router: `/?:uid`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

}
