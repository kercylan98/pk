package routers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context/param"
)

func init() {

    beego.GlobalControllerRouter["pk/controllers/olds:PlanController"] = append(beego.GlobalControllerRouter["pk/controllers/olds:PlanController"],
        beego.ControllerComments{
            Method: "Auto",
            Router: "/auto",
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["pk/controllers/olds:PlanController"] = append(beego.GlobalControllerRouter["pk/controllers/olds:PlanController"],
        beego.ControllerComments{
            Method: "Allows",
            Router: "/course/allows",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["pk/controllers/olds:PlanController"] = append(beego.GlobalControllerRouter["pk/controllers/olds:PlanController"],
        beego.ControllerComments{
            Method: "CourseMove",
            Router: "/course/move",
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["pk/controllers/olds:PlanController"] = append(beego.GlobalControllerRouter["pk/controllers/olds:PlanController"],
        beego.ControllerComments{
            Method: "Unallowable",
            Router: "/course/unallowable",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["pk/controllers/olds:PlanController"] = append(beego.GlobalControllerRouter["pk/controllers/olds:PlanController"],
        beego.ControllerComments{
            Method: "Draw",
            Router: "/draw",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["pk/controllers/olds:PlanController"] = append(beego.GlobalControllerRouter["pk/controllers/olds:PlanController"],
        beego.ControllerComments{
            Method: "Import",
            Router: "/import",
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["pk/controllers/olds:PlanController"] = append(beego.GlobalControllerRouter["pk/controllers/olds:PlanController"],
        beego.ControllerComments{
            Method: "NewPlan",
            Router: "/new",
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["pk/controllers/olds:PlanController"] = append(beego.GlobalControllerRouter["pk/controllers/olds:PlanController"],
        beego.ControllerComments{
            Method: "Optimize",
            Router: "/optimize",
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["pk/controllers/olds:PlanController"] = append(beego.GlobalControllerRouter["pk/controllers/olds:PlanController"],
        beego.ControllerComments{
            Method: "SectionMove",
            Router: "/section/move",
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["pk/controllers/olds:PlanController"] = append(beego.GlobalControllerRouter["pk/controllers/olds:PlanController"],
        beego.ControllerComments{
            Method: "SwitchOnline",
            Router: "/switch",
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["pk/controllers/olds:PlanController"] = append(beego.GlobalControllerRouter["pk/controllers/olds:PlanController"],
        beego.ControllerComments{
            Method: "Template",
            Router: "/template",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

}
