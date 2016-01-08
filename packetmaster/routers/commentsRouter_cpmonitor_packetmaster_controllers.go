package routers

import (
	"github.com/astaxie/beego"
)

func init() {
	
	beego.GlobalControllerRouter["cpmonitor/packetmaster/controllers:MasterController"] = append(beego.GlobalControllerRouter["cpmonitor/packetmaster/controllers:MasterController"],
		beego.ControllerComments{
			"Get",
			`/assapp/:srcendpoint`,
			[]string{"get"},
			nil})

}
