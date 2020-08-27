package controllers

import (
	"bailun.com/CT4_quote_server/WebManageSvr/models"
	"encoding/json"
	"github.com/astaxie/beego"
)

type FuncManagerController struct {
	beego.Controller
}

type dataBaseMapTable struct {
	dataBase string
	table    string
}

type ReturnGetTableConfig struct {
	models.CommonReturn `json:",inline"`
	Data                []models.DataTableConfig
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

type UpdateTBconfParm struct {
	FuncName string                         `json:"funcName"`
	Data     []models.DataTableUpdateConfig `json:"Data"`
}

// @Title update
// @Description update object
// @Param   body		body 	controllers.UpdateTBconfParm	true	     "parm"
// @Success 200
// @router /updateTableConfig [post]
func (f *FuncManagerController) UpdateTableConfig() {

	var ReturnData models.CommonReturn
	ReturnData.SetData(0, "success to update")
	defer func() {
		f.Data["json"] = ReturnData
		f.ServeJSON()

	}()
	var parmData UpdateTBconfParm
	err := json.Unmarshal(f.Ctx.Input.RequestBody, &parmData)
	if err != nil {
		ReturnData.SetData(1, "parm is error:"+err.Error())
		return
	}

	db, tb, ok := models.GetDBTBInfo(parmData.FuncName)
	if !ok {
		ReturnData.SetData(1, "not found table")
		return
	}
	err = models.UpdateDBConfig(db, tb, parmData.Data)
	if err != nil {
		ReturnData.SetData(1, err.Error())
		return
	}

}

type TBDataListReturn struct {
	models.CommonReturn `json:",inline"`
	Data                []map[string]interface{}
}
