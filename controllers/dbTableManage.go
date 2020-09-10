package controllers

import (
	"bailun.com/CT4_quote_server/WebManageSvr/models"
	"bailun.com/CT4_quote_server/WebManageSvr/service"
	"encoding/json"
	"github.com/astaxie/beego"
)

type DBTBInfoManagerController struct {
	beego.Controller
}

type ReturnGetAllTable struct {
	CommonReturn `json:",inline"`
	Data         []models.DBTBInfo
}

// @Title tableList
// @Description
// @Success 200  {object} controllers.ReturnGetAllTable
// @router /getAllTables [get]
func (f *DBTBInfoManagerController) GetAllTable() {
	//funcName := f.GetString("db", "")
	var ReturnData ReturnGetAllTable
	ReturnData.SetData(0, "success to get")
	defer func() {
		f.Data["json"] = ReturnData
		f.ServeJSON()
	}()

	dbTb, err := service.GetDBNames()
	if err != nil {
		ReturnData.SetData(1, err.Error())
		return
	}
	ReturnData.Data = dbTb
}

type UpdateTBConfParm struct {
	DB       string                         `json:"DB"`
	TB       string                         `json:"TB"`
	FuncName string                         `json:"funcName"`
	Data     []models.DataTableUpdateConfig `json:"Data"`
}

// @Title update
// @Description update object
// @Param   body		body 	controllers.UpdateTBConfParm	true	     "parm"
// @Success 200
// @router /updateTableConfig [post]
func (f *DBTBInfoManagerController) UpdateTableConfig() {

	var ReturnData CommonReturn
	ReturnData.SetData(0, "success to update")
	defer func() {
		f.Data["json"] = ReturnData
		f.ServeJSON()

	}()
	var parmData UpdateTBConfParm
	err := json.Unmarshal(f.Ctx.Input.RequestBody, &parmData)
	if err != nil {
		ReturnData.SetData(1, "parm is error:"+err.Error())
		return
	}

	err = service.UpdateDBConfig(parmData.DB, parmData.TB, parmData.Data)
	if err != nil {
		ReturnData.SetData(1, err.Error())
		return
	}

}

type ReturnGetTableConfig struct {
	CommonReturn `json:",inline"`
	Data         []models.DataTableConfigReturn
}

// @Title FuncList
// @Description list object
// @Param   DB  query   string     false       "db"
// @Param   TB  query   string     false       "tb"
// @Success  200  {object} controllers.ReturnGetTableConfig
// @router /TableConfig [get]
func (f *DBTBInfoManagerController) TableConfig() {
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
	res, err := service.GetDBTbConfig(&models.DbTb{db, tb})

	if err != nil {
		ReturnData.SetData(1, "get table config failed")
		return
	}

	ReturnData.Data = res
}
