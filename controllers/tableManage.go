package controllers

import (
	"bailun.com/CT4_quote_server/WebManageSvr/models"
	"github.com/astaxie/beego"
)

type TableManagerController struct {
	beego.Controller
}

// @Title FuncList
// @Description list object
// @Param   DB  query   string     false       "db"
// @Param   TB  query   string     false       "funcName"
// @Success 200
// @router /TableConfig [get]
func (f *TableManagerController) TableConfig() {
	db := f.GetString("DB", "")
	tb := f.GetString("TB", "")

	var ReturnData ReturnGetTableConfig
	ReturnData.SetData(0, "success to get")
	defer func() {
		f.Data["json"] = ReturnData
		f.ServeJSON()
	}()

	if tb == "" || db == "" {
		ReturnData.SetData(1, "not found table")
		return
	}
	res, err := models.GetDBTbConfig(db, tb)
	if err != nil {
		ReturnData.SetData(1, "get table config failed")
		return
	}
	ReturnData.Data = res
}

// @Title getTableDataList
// @Description getTableDataList
// @Param   funcName  query   string     false       "funcName"
// @Param   page  query   int     false       "funcName"
// @Param   size  query   int     false       "funcName"
// @Success 200
// @router /list [get]
func (f *TableManagerController) DataList() {
	var err error
	var page, size int

	db := f.GetString("DB", "")
	tb := f.GetString("TB", "")
	page, err = f.GetInt("page", 1)
	size, err = f.GetInt("size", 10)
	var ReturnData TBDataListReturn
	ReturnData.SetData(0, "success to get")
	defer func() {
		f.Data["json"] = ReturnData
		f.ServeJSON()

	}()

	var getData []map[string]interface{}
	getData, err = models.GetTableDataList(db, tb, page, size)
	if err != nil {
		ReturnData.SetData(1, err.Error())
	}

	ReturnData.Data = getData
	return
}
