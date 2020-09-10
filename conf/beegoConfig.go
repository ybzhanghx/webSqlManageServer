package conf

import (
	commConf "bailun.com/CT4_quote_server/common/conf"
	"encoding/json"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/plugins/cors"

	"net/http"
)

func BeeConfInit() {
	var logCfg = commConf.LoggerConfig{
		FileName:            Conf.LogConf.LogPath,
		Level:               Conf.LogConf.Level,
		EnableFuncCallDepth: true,
		LogFuncCallDepth:    3,
		RotatePerm:          "777",
		Perm:                "777",
		Rotate:              true,
		Maxsize:             1 << 28,
		Maxlines:            1000000,
	}

	// 设置beego log库的配置
	b, _ := json.Marshal(&logCfg)
	_ = logs.SetLogger(Conf.LogConf.Adapter, string(b))
	logs.EnableFuncCallDepth(logCfg.EnableFuncCallDepth)
	logs.SetLogFuncCallDepth(logCfg.LogFuncCallDepth)

	beego.BConfig.AppName = Conf.Server.AppName
	beego.BConfig.Listen.HTTPPort = Conf.Server.HTTPPort
	//beego.BConfig.RunMode = Conf.Server.RunMode
	// 是否允许在 HTTP 请求时，返回原始请求体数据字节，默认为 false
	beego.BConfig.CopyRequestBody = Conf.Server.CopyRequestBody
	// 是否模板自动渲染，默认值为 true，对于 API 类型的应用，应用需要把该选项设置为 false，不需要渲染模板
	beego.BConfig.WebConfig.AutoRender = Conf.Server.AutoRender
	beego.BConfig.WebConfig.EnableDocs = Conf.Server.EnableDocs
	beego.BConfig.RunMode = Conf.Server.RunMode
	beego.InsertFilter("*", beego.BeforeRouter, cors.Allow(&cors.Options{
		AllowAllOrigins:  true,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Authorization", "Access-Control-Allow-Origin", "Access-Control-Allow-Headers", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length", "Access-Control-Allow-Origin", "Access-Control-Allow-Headers", "Content-Type"},
		AllowCredentials: true,
	}))
	beego.ErrorHandler("404", func(rw http.ResponseWriter, r *http.Request) {
		_, _ = rw.Write([]byte("path not found"))
	})

	if beego.BConfig.RunMode == "dev" {
		beego.BConfig.WebConfig.DirectoryIndex = true
		beego.BConfig.WebConfig.StaticDir["/swagger"] = "swagger"
	}
}
