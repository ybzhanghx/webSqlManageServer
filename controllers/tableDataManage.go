package controllers

import (
	"WebManageSvr/models"
	"WebManageSvr/service"
	"encoding/json"
	"github.com/astaxie/beego"
)

type (
	TableDataManagerController struct {
		beego.Controller
	}
	//数据返回结构
	TBDataListReturn struct {
		CommonReturn `json:",inline"`
		Fields       []models.FieldType
		Data         []interface{}
		Totals       int64
	}

	//更新请求参数
	UpdateTableParm struct {
		DB  string
		TB  string
		Add string
		Del []int
		Upd string
	}
)

// @Title getTableDataList
// @Description getTableDataList
// @Param  DB   query   string     false       "funcName"
// @Param  TB   query   string     false       "funcName"
// @Param   page  query   int     false       "funcName"
// @Param   size  query   int     false       "funcName"
// @Success 200
// @router /list [get]
func (f *TableDataManagerController) DataList() {
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
	var dbtb = &models.DbTb{db, tb}
	getData, fields, err = service.GetTableDataList(dbtb, page, size)
	if err != nil {
		ReturnData.SetData(1, err.Error())
		return
	}
	ReturnData.Totals, err = service.GetTableDataTotals(dbtb)
	if err != nil {
		ReturnData.SetData(1, err.Error())
		return
	}
	ReturnData.Data = getData
	ReturnData.Fields = fields
	return
}

// @Title del
// @Description update object
// @Param   body		body 	controllers.UpdateTableParm	true	     "parm"
// @Success 200
// @router /updateTableData [post]
func (f *TableDataManagerController) UpdateTableData() {

	var ReturnData CommonReturn
	ReturnData.SetData(0, "success to update")
	defer func() {
		f.Data["json"] = ReturnData
		f.ServeJSON()

	}()
	var parmData UpdateTableParm
	err := json.Unmarshal(f.Ctx.Input.RequestBody, &parmData)
	if err != nil {
		ReturnData.SetData(1, "parm is error:"+err.Error())
		return
	}

	err = service.WriteTable(&models.DbTb{DB: parmData.DB, TB: parmData.TB}, parmData.Add, parmData.Upd, parmData.Del)
	if err != nil {
		ReturnData.SetData(1, err.Error())
		return
	}

}
