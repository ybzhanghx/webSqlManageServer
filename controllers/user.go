package controllers

import (
	"bailun.com/CT4_quote_server/front_gateway/models"
	"github.com/astaxie/beego"
)

type UserController struct {
	beego.Controller
}

// @Title userList
// @Description list object
// @Param   page  query   int     false       "page"
// @Param   size   query   int     false       "size limit"
// @Success 200
// @router /list [get]
func (c *UserController) List() {
	type Return struct {
		User  []models.UserWithNickName
		total int
		Code  int
		Msg   string
	}
	var ReturnData Return
	page, _ := c.GetInt("page", 1)
	size, _ := c.GetInt("size", 10)
	users, err := models.GetUsersPage(page, size)

	if err != nil {
		ReturnData.Code, ReturnData.Msg = 1, err.Error()
	} else {
		ReturnData.Code, ReturnData.Msg, ReturnData.User = 0, "success to get", users
	}
	c.Data["json"] = ReturnData
	c.ServeJSON()

}
