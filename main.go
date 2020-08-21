package main

import (
	"bailun.com/CT4_quote_server/WebManageSvr/conf"
	_ "bailun.com/CT4_quote_server/WebManageSvr/routers"
	"flag"

	"github.com/astaxie/beego"
)

func main() {
	var err error
	flag.Parse()
	if err = conf.Init(); err != nil {
		panic(err)
	}
	conf.BeeConfInit()
	conf.MysqlInit()
	if beego.BConfig.RunMode == "dev" {
		beego.BConfig.WebConfig.DirectoryIndex = true
		beego.BConfig.WebConfig.StaticDir["/swagger"] = "swagger"
	}
	beego.Run()
}
