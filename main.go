package main

import (
	"bailun.com/CT4_quote_server/WebManageSvr/conf"
	"bailun.com/CT4_quote_server/WebManageSvr/mysqls"
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
	mysqls.MysqlInit()

	beego.Run()
}
