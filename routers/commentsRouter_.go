package routers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context/param"
)

func init() {

	beego.GlobalControllerRouter["bailun.com/CT4_quote_server/WebManageSvr/controllers:DBTBInfoManagerController"] = append(beego.GlobalControllerRouter["bailun.com/CT4_quote_server/WebManageSvr/controllers:DBTBInfoManagerController"],
		beego.ControllerComments{
			Method:           "TableConfig",
			Router:           `/TableConfig`,
			AllowHTTPMethods: []string{"get"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["bailun.com/CT4_quote_server/WebManageSvr/controllers:DBTBInfoManagerController"] = append(beego.GlobalControllerRouter["bailun.com/CT4_quote_server/WebManageSvr/controllers:DBTBInfoManagerController"],
		beego.ControllerComments{
			Method:           "ConnectMysql",
			Router:           `/connectMysql`,
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["bailun.com/CT4_quote_server/WebManageSvr/controllers:DBTBInfoManagerController"] = append(beego.GlobalControllerRouter["bailun.com/CT4_quote_server/WebManageSvr/controllers:DBTBInfoManagerController"],
		beego.ControllerComments{
			Method:           "GetAllTable",
			Router:           `/getAllTables`,
			AllowHTTPMethods: []string{"get"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["bailun.com/CT4_quote_server/WebManageSvr/controllers:DBTBInfoManagerController"] = append(beego.GlobalControllerRouter["bailun.com/CT4_quote_server/WebManageSvr/controllers:DBTBInfoManagerController"],
		beego.ControllerComments{
			Method:           "UpdateTableConfig",
			Router:           `/updateTableConfig`,
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["bailun.com/CT4_quote_server/WebManageSvr/controllers:EntryManagerController"] = append(beego.GlobalControllerRouter["bailun.com/CT4_quote_server/WebManageSvr/controllers:EntryManagerController"],
		beego.ControllerComments{
			Method:           "FuncList",
			Router:           `/funcList`,
			AllowHTTPMethods: []string{"get"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["bailun.com/CT4_quote_server/WebManageSvr/controllers:EntryManagerController"] = append(beego.GlobalControllerRouter["bailun.com/CT4_quote_server/WebManageSvr/controllers:EntryManagerController"],
		beego.ControllerComments{
			Method:           "SaveFuncList",
			Router:           `/saveFuncList`,
			AllowHTTPMethods: []string{"get"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["bailun.com/CT4_quote_server/WebManageSvr/controllers:TableDataManagerController"] = append(beego.GlobalControllerRouter["bailun.com/CT4_quote_server/WebManageSvr/controllers:TableDataManagerController"],
		beego.ControllerComments{
			Method:           "DataList",
			Router:           `/list`,
			AllowHTTPMethods: []string{"get"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["bailun.com/CT4_quote_server/WebManageSvr/controllers:TableDataManagerController"] = append(beego.GlobalControllerRouter["bailun.com/CT4_quote_server/WebManageSvr/controllers:TableDataManagerController"],
		beego.ControllerComments{
			Method:           "UpdateTableData",
			Router:           `/updateTableData`,
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

}
