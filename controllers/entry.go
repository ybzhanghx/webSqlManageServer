package controllers

import (
	"WebManageSvr/conf"
	"WebManageSvr/service"
	"encoding/json"
	"github.com/astaxie/beego"
)

type (
	EntryManagerController struct {
		beego.Controller
	}
	ReturnFuncList struct {
		CommonReturn `json:",inline"`
		Data         conf.FuncListConf `json:"Data"`
	}
	ParmSaveFuncList struct {
		Data conf.FuncNode `json:"Data"`
	}
)

// @Title tableList
// @Description
// @Success 200  {object} controllers.ReturnFuncList
// @router /funcList [get]
func (f *EntryManagerController) FuncList() {
	var ReturnData ReturnFuncList
	ReturnData.SetData(0, successGetMsg)
	defer func() {
		f.Data["json"] = ReturnData
		f.ServeJSON()
	}()
	Data, err := service.GetFuncList()
	if err != nil {
		ReturnData.SetData(1, FailGetErr.Error())
		return
	}
	ReturnData.Data = *Data
}

// @Title tableList
// @Description
// @Param   body		body 	controllers.ParmSaveFuncList	true	     "parm"
// @Success 200  {object} controllers.ReturnFuncList
// @router /saveFuncList [get]
func (f *EntryManagerController) SaveFuncList() {
	var ReturnData CommonReturn
	ReturnData.SetData(0, successGetMsg)
	defer func() {
		f.Data["json"] = ReturnData
		f.ServeJSON()
	}()
	var parmData ParmSaveFuncList
	err := json.Unmarshal(f.Ctx.Input.RequestBody, &parmData)
	if err != nil {
		ReturnData.SetData(1, "parm is error:"+err.Error())
		return
	}

}
