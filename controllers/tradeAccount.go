package controllers

import (
	"bailun.com/CT4_quote_server/front_gateway/models"
	"bailun.com/CT4_quote_server/protocol/external"
	cpc "bailun.com/CT4_quote_server/protocol/rpc"
	"context"
	"encoding/json"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/jinzhu/copier"
)

type TradeAccountController struct {
	beego.Controller
}

// @Title userList
// @Description list object
// @Param   userId    query   string  false
// @Success 200
// @router /list [get]
func (this *TradeAccountController) List() {
	var returnData models.TradeAccountReturn
	returnData.SetData(0, "success to get")
	defer func() {
		this.Data["json"] = returnData
		this.ServeJSON()
	}()
	userId := this.GetString("userId", "")

	rpcReqData := new(cpc.TradeAccountEnquiryReqData)
	rpcReqData.Userid = userId

	rpcRes, err := rpcClient.TradeAccountListAll(context.Background(), &cpc.TradeAccountAllInfoReq{})
	if err != nil {
		logs.Error(err)
		returnData.SetData(1, "failed to get")
		return
	}

	returnData.Data = rpcRes.ResData

	return

}

// @Title update
// @Description update object
// @Param   body   body  models.TradAccountData  true  "用户数据"
// @Success 200  {object} models.CommonReturn
// @router /update [post]
func (this *TradeAccountController) Update() {
	var returnData models.CommonReturn
	returnData.SetData(0, "success to updated")
	defer func() {
		this.Data["json"] = returnData
		this.ServeJSON()
	}()
	ReqData := new(models.TradAccountData)
	err := json.Unmarshal(this.Ctx.Input.RequestBody, ReqData)

	if err != nil {
		returnData.SetData(1, "parm is failed")
		return
	}
	rpcReqData := new(cpc.TradeAccountItemInfo)
	_ = copier.Copy(rpcReqData, ReqData)
	_, err = rpcClient.TradeAccountUpdate(context.Background(), rpcReqData)
	if err != nil {
		logs.Error(err)
		returnData.SetData(1, "failed to update")
		return
	}
	//returnData.Data =*rpcRes

	return

}

// @Title Freeze
// @Description Freeze object
// @Param   body  body   models.TradeAccountSingReq true
// @Success 200 {object} models.CommonReturn
// @router /freeze [post]
func (this *TradeAccountController) Freeze() {
	var returnData models.CommonReturn
	returnData.SetData(0, "success to get")
	defer func() {
		this.Data["json"] = returnData
		this.ServeJSON()
	}()

	rpcReqData := new(models.TradeAccountSingReq)
	err := json.Unmarshal(this.Ctx.Input.RequestBody, rpcReqData)

	if err != nil {
		returnData.SetData(1, "failed to get")
		return
	}
	_, err = rpcClient.TradeAccountFreeze(context.Background(), &cpc.TradeAccountFreezeReq{
		ReqData: &cpc.TradeAccountFreezeReqData{AccountID: rpcReqData.AccountID},
	})
	if err != nil {
		logs.Error(err)
		returnData.SetData(1, err.Error())

		return
	}
	//returnData.Data =*rpcRes

	return

}

// @Title close
// @Description close object
// @Param   body  body   models.TradeAccountArrReq true
// @Success 200 {object} models.CommonReturn
// @router /close [post]
func (this *TradeAccountController) Close() {
	var returnData models.CommonReturn
	returnData.SetData(0, "success to close")
	defer func() {
		this.Data["json"] = returnData
		this.ServeJSON()
	}()

	httpReqData := new(models.TradeAccountArrReq)
	err := json.Unmarshal(this.Ctx.Input.RequestBody, httpReqData)

	if err != nil {
		returnData.SetData(1, "failed to close")
		return
	}

	_, err = rpcClient.TradeAccountClose(context.Background(),
		&cpc.TradeAccountCloseReq{ReqData: &external.TradeAccountsRepeated{
			AccountIDs: httpReqData.AccountIDs,
		}})
	if err != nil {
		logs.Error(err)
		returnData.SetData(1, err.Error())
		return
	}
	//returnData.SetData(0,)
	return

}
