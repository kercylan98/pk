package main

import (
	_ "pk/routers"
	//_ "pk/ibps/app"
	//_ "pk/ibps/engine/models"
	"github.com/astaxie/beego"
)

func slfadd(in int)(out int){
	out = in + 1
	return
}

func init() {
	beego.AddFuncMap("slfadd",slfadd)
}

func main() {
	//if beego.BConfig.RunMode == "dev" {
	//	beego.BConfig.WebConfig.DirectoryIndex = true
	//	beego.BConfig.WebConfig.StaticDir["/swagger"] = "swagger"
	//}
	// bee run -gendoc=true -downdoc=true
	beego.Run()
}

