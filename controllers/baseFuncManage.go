package controllers

import (
	"bailun.com/CT4_quote_server/front_gateway/models"
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

// @Title FuncList
// @Description list object
// @Param   funcName  query   string     false       "funcName"
// @Success 200
// @router /getTableConfig [get]
func (f *FuncManagerController) GetTableConfig() {
	funcName := f.GetString("funcName", "")
	var ReturnData ReturnGetTableConfig
	ReturnData.SetData(0, "success to get")
	defer func() {
		f.Data["json"] = ReturnData
		f.ServeJSON()
	}()
	db, tb, ok := models.GetDBTBInfo(funcName)
	if !ok {
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

type ReturnGetAllTable struct {
	models.CommonReturn `json:",inline"`
	Data                []string
}

// @Title tableList
// @Description list object
// @Param   db  query   string     false       "db"
// @Success 200
// @router /getAllTables [get]
func (f *FuncManagerController) GetAllTable() {
	funcName := f.GetString("db", "")
	var ReturnData ReturnGetAllTable
	ReturnData.SetData(0, "success to get")
	defer func() {
		f.Data["json"] = ReturnData
		f.ServeJSON()
	}()

	db, err := models.GetDBNames(funcName)
	if err != nil {
		ReturnData.SetData(1, err.Error())
		return
	}
	ReturnData.Data = db
}

type UpdateTBconfParm struct {
	FuncName string                         `json:"funcName"`
	Data     []models.DataTableUpdateConfig `json:"Data"`
}

// @Title FuncList
// @Description list object
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

// @Title getTableDataList
// @Description getTableDataList
// @Param   funcName  query   string     false       "funcName"
// @Param   page  query   int     false       "funcName"
// @Param   size  query   int     false       "funcName"
// @Success 200
// @router /getTableDataList [get]
func (f *FuncManagerController) GetTableDataList() {
	var err error
	var page, size int
	funcName := f.GetString("funcName", "")
	page, err = f.GetInt("page", 1)
	size, err = f.GetInt("size", 10)
	var ReturnData TBDataListReturn
	ReturnData.SetData(0, "success to get")
	defer func() {
		f.Data["json"] = ReturnData
		f.ServeJSON()

	}()
	db, tb, ok := models.GetDBTBInfo(funcName)
	if !ok {
		ReturnData.SetData(1, "failed to find tb")
	}

	var getData []map[string]interface{}
	getData, err = models.GetTableDataList(db, tb, page, size)
	if err != nil {
		ReturnData.SetData(1, err.Error())
	}

	ReturnData.Data = getData
	return
}
