package controllers

import (
	"bailun.com/CT4_quote_server/WebManageSvr/models"
	"github.com/astaxie/beego"
)

type FuncManagerController struct {
	beego.Controller
}

type ReturnGetAllTable struct {
	models.CommonReturn `json:",inline"`
	Data                []models.DBTBInfo
}

// @Title tableList
// @Description
// @Success 200  {object} controllers.ReturnGetAllTable
// @router /getAllTables [get]
func (f *FuncManagerController) GetAllTable() {
	//funcName := f.GetString("db", "")
	var ReturnData ReturnGetAllTable
	ReturnData.SetData(0, "success to get")
	defer func() {
		f.Data["json"] = ReturnData
		f.ServeJSON()
	}()

	dbTb, err := models.GetDBNames()
	if err != nil {
		ReturnData.SetData(1, err.Error())
		return
	}
	ReturnData.Data = dbTb
}
