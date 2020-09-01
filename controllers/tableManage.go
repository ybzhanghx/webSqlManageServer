package controllers

import (
	"bailun.com/CT4_quote_server/WebManageSvr/models"
	"encoding/json"
	"github.com/astaxie/beego"
)

type TableManagerController struct {
	beego.Controller
}

type ReturnGetTableConfig struct {
	models.CommonReturn `json:",inline"`
	Data                []models.DataTableConfig
}

// @Title FuncList
// @Description list object
// @Param   DB  query   string     false       "db"
// @Param   TB  query   string     false       "tb"
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

type TBDataListReturn struct {
	models.CommonReturn `json:",inline"`
	Fields              []models.FieldType
	Data                []interface{}
	Totals              int64
}

// @Title getTableDataList
// @Description getTableDataList
// @Param  DB   query   string     false       "funcName"
// @Param  TB   query   string     false       "funcName"
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

	var getData []interface{}
	var fields []models.FieldType
	getData, fields, err = models.GetTableDataList(db, tb, page, size)
	if err != nil {
		ReturnData.SetData(1, err.Error())
		return
	}
	ReturnData.Totals, err = models.GetTableDataTotals(db, tb)
	if err != nil {
		ReturnData.SetData(1, err.Error())
		return
	}
	ReturnData.Data = getData
	ReturnData.Fields = fields
	return
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
func (f *TableManagerController) UpdateTableConfig() {

	var ReturnData models.CommonReturn
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

	err = models.UpdateDBConfig(parmData.DB, parmData.TB, parmData.Data)
	if err != nil {
		ReturnData.SetData(1, err.Error())
		return
	}

}
