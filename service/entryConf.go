package service

import (
	"bailun.com/CT4_quote_server/WebManageSvr/conf"
	"bailun.com/CT4_quote_server/WebManageSvr/models"
	"github.com/astaxie/beego/logs"
)

func GetFuncList() (Res *conf.FuncListConf, err error) {
	if conf.Conf.FuncList.Types == "conf" {
		return conf.Conf.FuncList, nil
	}
	Res = &conf.FuncListConf{Name: "root", Value: "root", Children: []conf.FuncNode{}}
	var DBTBs []models.DBTBInfo
	DBTBs, err = GetDBNames()
	if err != nil {
		logs.Error(err.Error())
		return
	}
	for _, iv := range DBTBs {
		var tmp = conf.FuncNode{Name: iv.DbName, Value: iv.DbName, Children: []conf.FuncNode{}}
		for _, jv := range iv.TbName {
			tmp.Children = append(tmp.Children, conf.FuncNode{Value: iv.DbName + "|" + jv, Name: jv})
		}
		Res.Children = append(Res.Children, tmp)
	}
	return
}
