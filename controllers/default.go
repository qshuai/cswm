package controllers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/qshuai/cswm/models"
)

type MainController struct {
	beego.Controller
}

func (c *MainController) Get() {
	user := models.User{}
	user.Id = c.GetSession("uid").(int)

	o := orm.NewOrm()
	o.QueryTable("user").Filter("id", user.Id).One(&user, "ip", "last_login")

	c.Data["user"] = user
	c.Layout = "common.tpl"
	c.TplName = "index.html"
}
