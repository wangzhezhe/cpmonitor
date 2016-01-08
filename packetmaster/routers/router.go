package routers

import (
	"cpmonitor/packetmaster/controllers"
	"github.com/astaxie/beego"
)

func init() {
	beego.Include(&controllers.MasterController{})
	beego.SetStaticPath("/static", "static")
	beego.SetStaticPath("/doc", "static/doc")
}
