package main

import (
	_ "cpmonitor/packetmaster/routers"
	"github.com/astaxie/beego"
)

func main() {
	beego.Run()
}
